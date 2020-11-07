package user

import (
	"time"
	"math/rand"
)

type User struct {
	Id uint64
	LastTime time.Time
	IsAlive bool
}

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

func NewUser()*User{
	return &User{
		Id:r.Uint64(),
		LastTime: time.Now(),
		IsAlive:true,
	}
}

func (u *User)Alive(){
	u.LastTime = time.Now()
}

