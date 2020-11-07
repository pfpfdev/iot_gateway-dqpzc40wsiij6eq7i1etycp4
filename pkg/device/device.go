package device

import (
	"net"
	"time"
	"strings"
	. "../operable"
	"log"
	"fmt"
)

const defaultLifetime = 3

type Device struct{
	Name string
	conn net.Conn 
	Operables map[string]*Operable
	state *State
	lifetime uint32
	timer *time.Timer
	closeSig chan struct{}
	aliveSig chan struct{}
}

func NewDevice(name string,_conn net.Conn) *Device{
	return &Device{
		Name:name,
		conn:_conn,
		Operables:make(map[string]*Operable),
		state: NewState(),
		lifetime: defaultLifetime,
		closeSig:make(chan struct{}),
		aliveSig:make(chan struct{}),
	}
}

func (d *Device)addOperable(name string) (*Operable,error){
	d.Operables[name]=NewOperable(name)
	fmt.Printf("%#v\n",d)
	return d.Operables[name],nil
}

func (d *Device)finishInit()error{
	d.state.Set(CONNECTED)
	d.beginAliveChecker()
	return nil
}

func (d *Device)beginAliveChecker(){
	d.timer = time.NewTimer(time.Duration(d.lifetime) * time.Second)
	go func(){
		for{
			select{
			case <- d.timer.C:
				d.state.Set(DISCONNECTED)
			case <- d.closeSig:
				return 
			case <- d.aliveSig:
				if !d.timer.Stop() {
					<-d.timer.C
				}
				d.timer.Reset(time.Duration(d.lifetime) * time.Second)
			}
		}
	}()
}

func (d* Device)alive(){
	d.aliveSig <-struct{}{}
}

func (d* Device)Parse(data string){
	strs := strings.Split(data," ")
	if d.state.IsSame(INITIALIZING) {
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
			log.Print("[ REG ] Command", _cmd,"(",_type,")","was added on ",operableName," of ",d.Name)
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
	}else{
		log.Print("[ERROR] Wrong state of ", d.Name)
	}
}