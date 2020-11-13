package user

import (
	"time"
	"math/rand"
)

type User struct {
	id uint64
	LastTime time.Time
	Until time.Time
	IsAlive bool
}

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

func NewUser()*User{
	return &User{
		id:r.Uint64(),
		LastTime: time.Now(),
		IsAlive:true,
	}
}

func (u *User)Alive(){
	u.LastTime = time.Now()
}

type UserQueue []*User

func NewUserQueue()UserQueue{
	return make(UserQueue,0)
}

func (u *UserQueue)AddWaiting()uint64{
	*u = append(*u,NewUser())
	return (*u)[len(*u)-1].id
}

func (u *UserQueue)Alive(token uint64){
	for _,user := range *u{
		if user.id == token{
			user.Alive()
			return
		}
	}
}

func (u *UserQueue)Order(token uint64)int{
	for i := range *u{
		if (*u)[i].id == token{
			return i+1
		}
	}
	return -1
}

func (u *UserQueue)IsFront(token uint64)bool{
	return len(*u)>0&&(*u)[0].id == token
}

func (u *UserQueue)GetFront()uint64{
	if len(*u)>0{
		return (*u)[0].id
	}
	return 0
}

func (u *UserQueue)Next(lastAssign time.Time){
	count := len(*u)
	if count == 0{
		return
	}
	first := true
	*u = (*u)[1:]
	for i:= range *u{
		if (*u)[i].LastTime.After(lastAssign) {
			(*u)[i].IsAlive = true
		}else{
			(*u)[i].IsAlive = false
		}
		if first&&(*u)[i].IsAlive{
			*u = (*u)[i:]
			first = false
		}
	}
}
