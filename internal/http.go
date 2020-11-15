package internal

import (
	"net/http"
	"github.com/gorilla/mux"
	"log"
	"strconv"
)

func HttpServer(opt *HttpOpt){
	auth := NewBasicAuthMiddleware(opt)
	//Router機能はmuxを使用
	r:= mux.NewRouter()
	//数が少ないので一覧実装
	r.HandleFunc("/devices",DeviceList)
	r.HandleFunc("/devices/{name}",DeviceDetail)
	r.HandleFunc("/units",UnitList).Methods("GET")
	r.HandleFunc("/units",auth(MakeUnit)).Methods("POST")
	r.HandleFunc("/units/{name}",UnitDetail).Methods("GET")
	r.HandleFunc("/units/{name}",MakeBooking).Methods("POST")
	r.HandleFunc("/units/{name}/{operable}",Operate)
	r.HandleFunc("/log",auth(LogFetch))
	//ログを使用するように設定
	r.Use(loggingMiddleware)
	r.Use(corsMiddleware)
	//サーバーを設定して開始
	http.Handle("/",r)
	log.Print("Http Server Started on "+strconv.Itoa(*opt.Port))
	http.ListenAndServe(":"+strconv.Itoa(*opt.Port), nil)
}

func loggingMiddleware(next http.Handler) http.Handler {
	//ログを表示するためのミドルウェア
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Print("[HTTP] ",r.URL.Path," (",r.Method,") from ",r.RemoteAddr)
        next.ServeHTTP(w, r)
    })
}

func corsMiddleware(next http.Handler)http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        next.ServeHTTP(w, r)
    })
}

func NewBasicAuthMiddleware(opt *HttpOpt)func(fn http.HandlerFunc)http.HandlerFunc{
	basicUser := *opt.BasicAuth.User
	basicPass := *opt.BasicAuth.Password
	return func(fn http.HandlerFunc)http.HandlerFunc{
		return basicAuthMiddleware(fn,basicUser,basicPass)
	}
}

func basicAuthMiddleware(fn http.HandlerFunc, basicUser string, basicPass string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, pass, _ := r.BasicAuth()
		if user != basicUser || pass != basicPass {
			http.Error(w, "Unauthorized.", 401)
			return
		}
		fn(w, r)
	}
}