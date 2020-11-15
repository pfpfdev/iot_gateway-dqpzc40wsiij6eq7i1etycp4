package internal

import (
	"time"
	"log"
)

func DeviceGC(opt *StrategyOpt){
	//定期的に使用されていないデバイスを削除する
	cycle := time.Duration(*opt.DeviceCycle) * time.Second
	for{
		for name,device := range Devices {
			if device.LastAlive.Before(time.Now().Add(-cycle)) {
				log.Print(name," was deleted by GC")
				device.Close()
				delete(Devices,name)
			}
		}
		time.Sleep(cycle)
	}
}