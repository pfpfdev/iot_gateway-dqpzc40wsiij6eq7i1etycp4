package unit

import (
	"sync"
	"time"
	. "../user"
	. "../device"
	"log"
)

const cycle = time.Second * 90

type Unit struct {
	Name string
	Operables map[string]*Operable
	Queue UserQueue
	ticker *time.Ticker
	lastAssign time.Time
	mutex sync.Mutex
}

func NewUnit(name string)*Unit{
	log.Print("New Unit ", name," Created")
	return &Unit{
		Name:name,
		Operables: make(map[string]*Operable),
		Queue:NewUserQueue(),
		lastAssign:time.Now(),
	}
}

func (u *Unit)Book()uint64{
	log.Print("New user made booking of ",u.Name)
	u.mutex.Lock()
	token := u.Queue.AddWaiting()
	if len(u.Queue) == 1 {
		log.Print("User(",u.Queue.GetFront(),") get controls of ",u.Name)
		u.Queue[0].Until = time.Now().Add(cycle)
		u.Start()
	}
	u.mutex.Unlock()
	return token
}

func (u *Unit)Start(){
	go func(){
		u.ticker = time.NewTicker(cycle)
		for{
			<- u.ticker.C
			u.mutex.Lock()
			u.Queue.Next(u.lastAssign)
			log.Print(len(u.Queue)," Users are waiting for ",u.Name)
			if len(u.Queue) == 0{
				break
			}else{
				log.Print("User(",u.Queue.GetFront(),") get controls of ",u.Name)
				u.Queue[0].Until = time.Now().Add(cycle)
			}
			u.lastAssign = time.Now()
			u.mutex.Unlock()
		}
	}()
}