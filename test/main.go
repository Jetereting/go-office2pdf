package main

import (
	"fmt"
	"github.com/astaxie/beego/httplib"
	"time"
)

var channel = make(chan int, 5)

func main() {
	channel <- 112
	channel <- 112
	channel <- 112
	channel <- 112
	channel <- 112
	go func() {
		<-channel
		fmt.Println("Hello1")
		httplib.Post("http://127.0.0.1:3000/convert").PostFile("file", "/Users/au/Desktop/软著/黄金价格提醒工具说明书.doc").ToFile("1.pdf")
		channel <- 1
	}()
	go func() {
		<-channel
		fmt.Println("Hello2")
		httplib.Post("http://127.0.0.1:3000/convert").PostFile("file", "/Users/au/Desktop/软著/黄金价格提醒工具说明书.doc").ToFile("2.pdf")
		channel <- 1
	}()
	go func() {
		<-channel
		fmt.Println("Hello3")
		httplib.Post("http://127.0.0.1:3000/convert").PostFile("file", "/Users/au/Desktop/软著/黄金价格提醒工具说明书.doc").ToFile("3.pdf")
		channel <- 1
	}()
	go func() {
		<-channel
		fmt.Println("Hello4")
		httplib.Post("http://127.0.0.1:3000/convert").PostFile("file", "/Users/au/Desktop/软著/黄金价格提醒工具说明书.doc").ToFile("4.pdf")
		channel <- 1
	}()
	go func() {
		<-channel
		fmt.Println("Hello5")
		httplib.Post("http://127.0.0.1:3000/convert").PostFile("file", "/Users/au/Desktop/软著/黄金价格提醒工具说明书.doc").ToFile("5.pdf")
		channel <- 1
	}()
	go func() {
		<-channel
		fmt.Println("Hello6")
		httplib.Post("http://127.0.0.1:3000/convert").PostFile("file", "/Users/au/Desktop/软著/黄金价格提醒工具说明书.doc").ToFile("6.pdf")
		channel <- 1
	}()
	time.Sleep(time.Second * 20)
}
