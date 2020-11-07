package internal

import (
	"encoding/json"
	"../pkg/device"
	"net/http"
	"fmt"
	"github.com/gorilla/mux"
)

var Devices map[string]*device.Device = make(map[string]*device.Device)



func DeviceList(w http.ResponseWriter, r *http.Request){
	data,_ := json.Marshal(Devices)
	fmt.Fprintln(w,string(data))
}

func DeviceDetail(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	name := vars["name"]
	if _,ok := Devices[name];!ok{
		errMsg,_ := json.Marshal(map[string]interface{}{"Error":"Undefined Device"})
		fmt.Fprintln(w,string(errMsg))
		return
	}
	data,_ := json.Marshal(Devices[name])
	fmt.Fprintln(w,string(data))
}
