package main

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

var channel = make(chan int, 10)

func main() {
	channel <- 1
	log.Println("Server listen port: 3000")
	http.HandleFunc("/convert", handler)
	_ = http.ListenAndServe(":3000", nil)
}

func downloadFile(path string, url string) error {

	// Create the file
	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	keys, ok := r.URL.Query()["fileSrc"]

	if !ok || len(keys[0]) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("fileSrc Param 'key' is missing"))
		return
	}

	fileSrc := keys[0]

	log.Println("request incoming with fileSrc:", fileSrc)

	workDir := os.TempDir()

	fileName := strings.TrimSuffix(filepath.Base(fileSrc), filepath.Ext(fileSrc)) + time.Now().Format("060102150405")

	inputFile := filepath.Join(workDir, fileName)

	outputFile := filepath.Join(workDir, fileName+".pdf")

	errorDownload := downloadFile(inputFile, fileSrc)
	log.Println("下载耗时(秒):", time.Now().Sub(startTime).Seconds())

	if errorDownload != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("error when download file"))
		return
	}

	<-channel
	cmd := exec.Command("libreoffice", "--headless", "--convert-to", "pdf:writer_pdf_Export", inputFile, "--outdir", workDir)
	_, err := cmd.Output()
	channel <- 1

	if err != nil {
		log.Println("cmd err:", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename="+filepath.Base(outputFile))
	w.WriteHeader(http.StatusOK)

	fileBytes, err := ioutil.ReadFile(outputFile)
	if err != nil {
		log.Println("read file err:", err.Error())
		return
	}

	_, err = w.Write(fileBytes)
	if err != nil {
		log.Println("write file err:", err.Error())
		return
	}

	_ = os.Remove(outputFile)
	_ = os.Remove(inputFile)

	log.Println("总耗时(秒):", time.Now().Sub(startTime).Seconds())

	return
}
