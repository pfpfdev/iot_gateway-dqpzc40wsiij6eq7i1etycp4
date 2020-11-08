package internal

import (
	"net"
	"../pkg/device"
	"log"
	"time"
)

func SocketServer(){
	//サーバーの開始
	listener, err := net.Listen("tcp", "0.0.0.0:8081")

    if err != nil {
        panic(err)
	}
	log.Print("Socket Server Started")
	defer listener.Close()

	go DeviceGC()

	//コネクションのハンドリング
    for {
		//接続要求を処理する
		conn, err := listener.Accept()
		defer conn.Close()

        if err != nil {
			log.Print("Failed to Establish the Socket Connection ", err.Error())
            continue
		}
		log.Print("New Socket Connection Established From ",conn.RemoteAddr())

		//接続が確率したら別スレッドで処理を開始する
        go handler(conn)

    }
}

func handler(conn net.Conn){
	//デバイスのインスタンスを作成
	d:=device.NewDevice(conn)
	go d.Communicate()
	d.WaitEvent()
	//管理対象に追加
	Devices[d.Name]=d
	log.Print("New Device Defined : ",d.Name,"(",conn.RemoteAddr(),")")
	d.WaitEvent()
}

func DeviceGC(){
	const Cycle = 10 * time.Second
	for{
		for name,device := range Devices {
			if device.LastAlive.Before(time.Now().Add(-Cycle)) {
				log.Print(name," was deleted by GC")
				device.Close()
				delete(Devices,name)
			}
		}
		time.Sleep(Cycle)
	}
}