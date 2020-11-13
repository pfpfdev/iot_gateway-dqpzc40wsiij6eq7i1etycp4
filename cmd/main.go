package main

import (
	"../internal"
	"log"
	"os"
	"io"
)

func main(){
	//エントリポイント
	//ログファイルを開きながらサーバーを二つ起動する
	f,err:= os.OpenFile("/tmp/iot_gateway.log",os.O_APPEND|os.O_CREATE|os.O_WRONLY,0666)
	if err != nil {
		log.Fatal("Failed to set the log file")
	}
	defer f.Close()
	log.SetOutput(io.MultiWriter(f,os.Stdout))

	go internal.HttpServer()
	internal.SocketServer()
}