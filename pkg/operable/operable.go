package operable

import (
	"bufio"
	"fmt"
)

type Operable struct{
	Name string
	Operations map[string]Operation
	device *bufio.ReadWriter
}

func NewOperable(name string) *Operable{
	return &Operable{
		Name:name,
		Operations:make(map[string]Operation),
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
	val,exist := o.Operations[_cmd]
	if !exist{
		return "",UndefinedOperationErr()
	}
	if !definedType[val.Type].Match([]byte(_arg)){
		return "",InvalidArgumentErr()
	}
	fmt.Fprintf(o.device,"%s,%s,%s\n",o.Name,_cmd,_arg)
	res,_,err :=o.device.ReadLine()
	if err!=nil{
		return "",err
	}
	return string(res),nil
}