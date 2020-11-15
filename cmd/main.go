package main

import (
	"../internal"
	"log"
	"os"
	"io"
	"../pkg/unit"
)

func main(){
	config,err := internal.ParseYaml("./config.yaml")
	if err != nil{
		log.Fatal(err.Error())
	}
	//エントリポイント
	//ログファイルを開きながらサーバーを二つ起動する
	f,err:= os.OpenFile(*config.LogPath,os.O_APPEND|os.O_CREATE|os.O_WRONLY,0666)
	if err != nil {
		log.Fatal("Failed to set the log file")
	}
	defer f.Close()

	log.SetOutput(io.MultiWriter(f,os.Stdout))
	unit.SetCycle(*config.Strategy.ControlCycle)

	go internal.HttpServer(config.HttpServer)
	go internal.SocketServer(config.SocketServer)
	internal.DeviceGC(config.Strategy)
}