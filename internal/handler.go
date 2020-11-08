package internal

import (
	"encoding/json"
	"../pkg/device"
	"../pkg/unit"
	"../pkg/user"
	"net/http"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"strconv"
)

var Devices map[string]*device.Device = make(map[string]*device.Device)
var Units map[string]*unit.Unit = make(map[string]*unit.Unit)


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

func UnitList(w http.ResponseWriter, r *http.Request){
	data,_ := json.Marshal(Units)
	fmt.Fprintln(w,string(data))
}

func UnitDetail(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	name := vars["name"]
	if _,ok := Units[name];!ok{
		errMsg,_ := json.Marshal(map[string]interface{}{"Error":"Undefined Unit"})
		fmt.Fprintln(w,string(errMsg))
		return
	}
	token := r.URL.Query().Get("token")
	if len(token) != 0 {
		token ,err:= strconv.ParseUint(token,10,64)
		if err!=nil{
			errMsg,_ := json.Marshal(map[string]interface{}{"Error":"The token is invalid format"})
			fmt.Fprintln(w,string(errMsg))
			return
		}
		for _,u := range Units[name].Queue{
			if u.Id == token{
				u.Alive()
				data,_ := json.Marshal(Units[name])
				fmt.Fprintln(w,string(data))
				return
			}
		}
		errMsg,_ := json.Marshal(map[string]interface{}{"Error":"The token wasn't issued"})
		fmt.Fprintln(w,string(errMsg))
		return
	}else{
		data,_ := json.Marshal(Units[name])
		fmt.Fprintln(w,string(data))
	}
}

func MakeUnit(w http.ResponseWriter, r *http.Request){
	var data interface{}
	body, _ := ioutil.ReadAll(r.Body)
	println(string(body))
	err := json.Unmarshal(body,&data)
	if err!=nil{
		println("[ERROR]",err.Error())
	}
	fmt.Printf("%#v\n",data)
	units,ok := data.(map[string]interface{})
	if !ok{
		return
	}
	for unitName,devices := range units {
		devices,ok := devices.(map[string]interface{})
		if !ok {
			continue
		}
		Units[unitName]=unit.NewUnit(unitName)
		for deviceName,operables := range devices{
			device, check1 := Devices[deviceName]
			operables,check2 := operables.([]interface{})
			if !check1 || !check2 {
				continue
			}
			for _,operableName := range operables {
				operableName, check := operableName.(string)
				if !check{
					continue
				} 
				operable,ok := device.Operables[operableName]
				if !ok{
					continue
				}
				Units[unitName].Operables[operableName]=operable
			}
		}
	}
}

func MakeBooking(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	name := vars["name"]
	if _,ok := Units[name];!ok{
		errMsg,_ := json.Marshal(map[string]interface{}{"Error":"Undefined Unit"})
		fmt.Fprintln(w,string(errMsg))
		return
	}
	u:=user.NewUser() 
	Units[name].Book(u)
	fmt.Fprintln(w,map[string]interface{}{"token":u.Id})
}

func Operate(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	unitName := vars["name"]
	operableName := vars["operable"]
	unit,ok := Units[unitName]
	if !ok{
		errMsg,_ := json.Marshal(map[string]interface{}{"Error":"Undefined Unit"})
		fmt.Fprintln(w,string(errMsg))
		return
	}
	operable,ok := unit.Operables[operableName]
	if !ok{
		errMsg,_ := json.Marshal(map[string]interface{}{"Error":"Undefined Operable"})
		fmt.Fprintln(w,string(errMsg))
		return
	}
	tokenStr := r.URL.Query().Get("token")
	if len(tokenStr) == 0 {
		errMsg,_ := json.Marshal(map[string]interface{}{"Error":"The token is empty"})
		fmt.Fprintln(w,string(errMsg))
		return
	}
	token ,err:= strconv.ParseUint(tokenStr,10,64)
	if err!=nil{
		errMsg,_ := json.Marshal(map[string]interface{}{"Error":"The token is invalid format"})
		fmt.Fprintln(w,string(errMsg))
		return
	}
	cmd := r.URL.Query().Get("token")
	if len(cmd) != 0 {
		errMsg,_ := json.Marshal(map[string]interface{}{"Error":"The cmd is empty"})
		fmt.Fprintln(w,string(errMsg))
		return
	}
	if unit.User.Id == token {
		operable.Operate(r.URL.Query().Get("cmd"),r.URL.Query().Get("arg"))
	}
}