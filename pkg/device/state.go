package device

import (
	"sync"
)

const (
	INITIALIZING = iota
	CONNECTED
	DISCONNECTED
)

type State struct{
	value int
	mutex sync.Mutex
}

func NewState()*State{
	return &State{
		value:INITIALIZING,
	}
}

func (s *State)Set(_value int){
	s.mutex.Lock()
	s.value  =  _value
	s.mutex.Unlock()
}

func (s *State)IsSame(_value int)bool{
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return _value == s.value
}