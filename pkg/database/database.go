package database

import (
	"sync"
)

type Database struct {
	data map[string]map[string]interface{}
	mutex sync.Mutex
}

func NewDatabase() *Database{
	return &Database{
		data: make(map[string]map[string]interface{}),
	}
}

func (d *Database)NewCategory(name string){
	d.mutex.Lock()
	d.data[name] = make(map[string]interface{})
	d.mutex.Unlock()
}

func (d *Database)Insert(cat string, key string, value interface{}){
	subdata,exist := d.data[cat]
	if !exist{
		return 
	}
	d.mutex.Lock()
	subdata[key] = value
	d.mutex.Unlock()
}

func (d *Database)Get(cat string, key string)interface{}{
	subdata,exist := d.data[cat]
	if !exist{
		return nil
	}
	d.mutex.Lock()
	defer d.mutex.Unlock()
	value,exist := subdata[key]
	if !exist{
		return nil
	}
	return value
}