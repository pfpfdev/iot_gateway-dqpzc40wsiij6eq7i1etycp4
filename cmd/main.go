package main

import (
	"../internal"
)

func main(){
	go internal.HttpServer()
	internal.SocketServer()
}