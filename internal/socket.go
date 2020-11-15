package internal

import (
	"net"
	"../pkg/device"
	"log"
	"strconv"
)

func SocketServer(opt *SocketOpt){
	//サーバーの開始
	listener, err := net.Listen("tcp", "0.0.0.0:"+strconv.Itoa(*opt.Port))

    if err != nil {
        panic(err)
	}
	log.Print("Socket Server Started on "+strconv.Itoa(*opt.Port))

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