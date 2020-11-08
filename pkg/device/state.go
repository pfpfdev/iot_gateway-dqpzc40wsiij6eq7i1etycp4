package device

import (
	"sync"
)

const (
	INITIALIZING = "INITIALIZING"
	CONNECTED = "CONNECTED"
	DISCONNECTED = "DISCONNECTED"
)

type State struct{
	value string
	mutex sync.Mutex
}

func NewState()*State{
	return &State{
		value:INITIALIZING,
	}
}

func (s *State)Set(_value string){
	s.mutex.Lock()
	s.value  =  _value
	s.mutex.Unlock()
}

func (s *State)IsSame(_value string)bool{
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return _value == s.value
}