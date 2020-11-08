package device

import (
	"net"
	"time"
	"strings"
	. "../operable"
	"log"
	"fmt"
	"bufio"
)

type Device struct{
	Name string
	conn net.Conn 
	Operables map[string]*Operable
	state *State
	LastAlive time.Time
}

func NewDevice(name string,_conn net.Conn) *Device{
	return &Device{
		Name:name,
		conn:_conn,
		Operables:make(map[string]*Operable),
		state: NewState(),
		LastAlive:time.Now(),
	}
}

func (d *Device)addOperable(name string) (*Operable,error){
	w:= bufio.NewWriter(d.conn)
	r:= bufio.NewReader(d.conn)
	wr := bufio.NewReadWriter(r,w)
	d.Operables[name]=NewOperable(name,wr)
	fmt.Printf("%#v\n",d)
	return d.Operables[name],nil
}

func (d *Device)finishInit()error{
	d.state.Set(CONNECTED)
	return nil
}

func (d *Device)Close(){
	d.conn.Close()
}

func (d* Device)Parse(data string){
	d.LastAlive = time.Now()
	strs := strings.Split(data," ")
	if d.state.IsSame(INITIALIZING) {
		log.Print("String From ",d.Name," : ",data)
		switch strs[0] {
		case "ADD":
			if len(strs) !=2{
				log.Print("[ERROR] Log Argument")
				return
			}
			operableName := strs[1]
			//エラーは起きない
			d.addOperable(operableName)
			log.Print("[ ADD ] Operable ",operableName," was added to ",d.Name)
			return
		case "REG":
			if len(strs) !=4{
				log.Print("[ERROR] Log Argument")
				return
			}
			operableName := strs[1]
			_cmd := strs[2]
			_type:= strs[3]
			if _,ok := d.Operables[operableName];!ok{
				log.Print("[ERROR] Undefined operable")
				return
			}
			err := d.Operables[operableName].RegisterOperation(_cmd,_type)
			if err!= nil{
				log.Print("[ERROR] ",err.Error()," on ",d.Name)
				return
			}
			log.Print("[ REG ] Command ", _cmd,"(",_type,")"," was added on ",operableName," of ",d.Name)
			return
		case "FIN":
			if len(strs) !=1{
				log.Print("[ERROR] Log Argument")
				return
			}
			//エラーは起きない
			d.finishInit()
			log.Print("[ FIN ] Initialization fnished of ", d.Name)
			return
		default:
			log.Print("[ERROR] Undefined Commands")
		}
	}
}