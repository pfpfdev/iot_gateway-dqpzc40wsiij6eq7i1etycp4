package device

import (
	"fmt"
)

type Operable struct{
	Name string
	Operations map[string]Operation
	device *Device
}

func NewOperable(name string,_device *Device) *Operable{
	return &Operable{
		Name:name,
		Operations:make(map[string]Operation),
		device:_device,
	}
}

func (o *Operable)RegisterOperation(_cmd string,_type string)error{
	_,exist := definedType[_type]
	if !exist{
		return UndefinedTypeErr()
	}
	o.Operations[_cmd]=Operation{
		Cmd:_cmd,
		Type:_type,
	}
	return nil
}

func (o *Operable)Operate(_cmd string,_arg string)(string,error){
	println(_arg)
	val,exist := o.Operations[_cmd]
	if !exist{
		return "",UndefinedOperationErr()
	}
	if !(definedType[val.Type].Match([]byte(_arg))){
		return "",InvalidArgumentErr()
	}
	o.device.mutex.Lock()
	fmt.Fprintf(o.device,"%s %s %s\n",o.Name,_cmd,_arg)
	res := o.device.ReadLine()
	o.device.mutex.Unlock()
	return res,nil
}