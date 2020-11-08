package unit

import (
	"sync"
	"time"
	. "../user"
	. "../operable"
	"log"
)

const cycle = time.Second * 90

type Unit struct {
	Name string
	Operables map[string]*Operable
	Queue []*User
	ticker *time.Ticker
	User *User
	lastAssign time.Time
	mutex sync.Mutex
}

func NewUnit(name string)*Unit{
	log.Print("New Unit ", name," Created")
	return &Unit{
		Name:name,
		Operables: make(map[string]*Operable),
		Queue:make([]*User,0),
		User: nil,
		lastAssign:time.Now(),
	}
}

func (u *Unit)Book(user *User){
	log.Print("New user made booking of ",u.Name)
	u.mutex.Lock()
	if u.User == nil{
		u.SetUser(user)
	}else{
		u.Queue = append(u.Queue,user)
	}
	u.mutex.Unlock()
}

func (u *Unit)Assign(){
	count := len(u.Queue)
	log.Print("Assigning ",u.Name," : ", count," users are waiting")
	if count == 0{
		u.User = nil
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
	log.Print("User(",user.Id,") get controls of ",u.Name)
	u.User = user
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