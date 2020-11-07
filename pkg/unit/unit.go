package unit

import (
	"sync"
	"time"
	. "../user"
	. "../operable"
)

const cycle = time.Second * 90

type Unit struct {
	Name string
	Operables map[string]Operable
	Queue []*User
	ticker *time.Ticker
	user *User
	lastAssign time.Time
	mutex sync.Mutex
}

func NewUnit(name string)*Unit{
	return &Unit{
		Name:name,
		Operables: make(map[string]Operable),
		Queue:make([]*User,0),
		user: nil,
		lastAssign:time.Now(),
	}
}

func (u *Unit)Book(user *User){
	u.mutex.Lock()
	if u.user == nil{
		u.SetUser(user)
	}else{
		u.Queue = append(u.Queue,user)
	}
	u.mutex.Unlock()
}

func (u *Unit)Assign(){
	count := len(u.Queue)
	if count == 0{
		u.user = nil
		return 
	}
	first := true
	for i:= range u.Queue{
		if u.Queue[i].LastTime.After(u.lastAssign) {
			u.Queue[i].IsAlive = false
		}
		if first {
			u.SetUser(u.Queue[i])
			u.Queue = u.Queue[i+1:]
			first = false
		}
	}
}

func (u *Unit)SetUser(user *User){
	u.user = user
	u.ticker = time.NewTicker(cycle)
}

func (u *Unit)Start(){
	go func(){
		for{
			<- u.ticker.C
			u.mutex.Lock()
			u.Assign()
			u.mutex.Unlock()
		}
	}()
}