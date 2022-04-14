package main

import (
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

var channel = make(chan int, 1)

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
	keys, ok := r.URL.Query()["originFile"]

	isFileFlag := false
	originFile := ""
	var partFile *multipart.Part
	if !ok || len(keys[0]) < 1 {
		isFileFlag = true
	}
	if isFileFlag {
		reader, err := r.MultipartReader()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		part, err := reader.NextPart()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		partFile = part
		originFile = part.FileName()
	}
	if !isFileFlag {
		originFile = keys[0]
	}

	log.Println("request incoming with originFile:", originFile)

	workDir := os.TempDir()

	<-channel
	defer func() {
		channel <- 1
	}()
	newFileName := strings.TrimSuffix(filepath.Base(originFile), filepath.Ext(originFile)) + time.Now().Format("060102150405")
	serverFile := filepath.Join(workDir, newFileName)
	outputFile := filepath.Join(workDir, newFileName+".pdf")

	if isFileFlag {
		f, _ := os.Create(serverFile)
		_, _ = io.Copy(f, partFile)
	} else {
		// 如果是文件网址，则下载文件
		errorDownload := downloadFile(serverFile, originFile)
		if errorDownload != nil {
			log.Println("下载失败:", errorDownload)
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte("下载失败"))
			return
		}
		log.Println("下载耗时(秒):", time.Now().Sub(startTime).Seconds())
	}

	cmd := exec.Command("libreoffice", "--headless", "--convert-to", "pdf:writer_pdf_Export", serverFile, "--outdir", workDir)
	_, err := cmd.Output()
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
	_ = os.Remove(serverFile)

	log.Println("总耗时(秒):", time.Now().Sub(startTime).Seconds())

	return
}
