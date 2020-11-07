package internal

import (
	"net/http"
	"github.com/gorilla/mux"
)

func HttpServer(){
	r:= mux.NewRouter()
	r.HandleFunc("/devices",DeviceList)
	r.HandleFunc("/devices/{name}",DeviceDetail)
	http.Handle("/",r)
	http.ListenAndServe(":8080", nil)
}