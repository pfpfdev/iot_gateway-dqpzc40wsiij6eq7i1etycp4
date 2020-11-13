package device

import (
	"net"
	"time"
	"strings"
	"log"
	"bufio"
	"sync"
)

type Device struct{
	Name string
	conn net.Conn
	Operables map[string]*Operable
	state *State
	LastAlive time.Time
	mutex sync.Mutex
	buffer chan string
	event chan struct{}
}

func NewDevice(conn net.Conn) *Device{
	return &Device{
		conn:conn,
		Operables:make(map[string]*Operable),
		state: NewState(),
		LastAlive:time.Now(),
		buffer: make(chan string,1),
		event: make(chan struct{}),
	}
}

func (d *Device)addOperable(name string) (*Operable,error){
	d.Operables[name]=NewOperable(name,d)
	return d.Operables[name],nil
}

func (d *Device)finishInit()error{
	d.state.Set(CONNECTED)
	d.event <- struct{}{}
	return nil
}

func (d *Device)Close()error{
	close(d.event)
	close(d.buffer)
	return d.conn.Close()
}

func (d *Device)ReadLine()string{
	return <- d.buffer
}

func (d *Device)Write(p []byte)(n int,err error){
	return d.conn.Write(p)
}

func (d *Device)WaitEvent(){
	<- d.event
}

func (d *Device)Communicate(){
	r := bufio.NewReader(d.conn)
	name,err := r.ReadString('\n')
	if err!=nil{
		log.Print("Socket Connection ",d.conn.RemoteAddr()," made error ",err.Error())
		return
	}
	//CSV的に処理するのでカンマを消す
	name = strings.Replace(name," ","",-1)
	//最後に改行があるので削除する
	d.Name = name[:len(name)-1]
	for {
		//入力された文字をdeviceで処理する
		str,err := r.ReadString('\n')
		if err!=nil{
			log.Print("[ERROR] Failed to communicate ",err.Error())
			return
		}
		print(str)
		d.Parse(str[:len(str)-1])
	}
}

func (d* Device)Parse(data string){
	d.LastAlive = time.Now()
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
	}else{
		if len(data)!=0{
			println(data)
			d.buffer <- data
		}
	}
}