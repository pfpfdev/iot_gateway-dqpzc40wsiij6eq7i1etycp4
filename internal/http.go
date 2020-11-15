package internal

import (
	"net/http"
	"github.com/gorilla/mux"
	"log"
	"os"
)

func HttpServer(){
	//Router機能はmuxを使用
	r:= mux.NewRouter()
	//数が少ないので一覧実装
	r.HandleFunc("/devices",DeviceList)
	r.HandleFunc("/devices/{name}",DeviceDetail)
	r.HandleFunc("/units",UnitList).Methods("GET")
	r.HandleFunc("/units",ManageFunction(MakeUnit)).Methods("POST")
	r.HandleFunc("/units/{name}",UnitDetail).Methods("GET")
	r.HandleFunc("/units/{name}",MakeBooking).Methods("POST")
	r.HandleFunc("/units/{name}/{operable}",Operate)
	r.HandleFunc("/log",LogFetch)
	//ログを使用するように設定
	r.Use(loggingMiddleware)
	//サーバーを設定して開始
	http.Handle("/",r)
	http.ListenAndServe(":8080", nil)
}

func loggingMiddleware(next http.Handler) http.Handler {
	//ログを表示するためのミドルウェア
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Print("[HTTP] ",r.URL.Path," (",r.Method,") from ",r.RemoteAddr)
        next.ServeHTTP(w, r)
    })
}


var basicUser = "kokenuser"
var basicPass = os.Getenv("PASSWORD")

func ManageFunction(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, pass, _ := r.BasicAuth()
		if user != basicUser || pass != basicPass {
			http.Error(w, "Unauthorized.", 401)
			return
		}
		fn(w, r)
	}
}