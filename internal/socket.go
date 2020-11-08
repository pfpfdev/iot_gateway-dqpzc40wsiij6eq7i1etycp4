package internal

import (
	"net"
	"bufio"
	"../pkg/device"
	"log"
	"strings"
)

func SocketServer(){
	//サーバーの開始
	listener, err := net.Listen("tcp", "0.0.0.0:8081")

    if err != nil {
        panic(err)
	}
	log.Print("Socket Server Started")
	defer listener.Close()

	//コネクションのハンドリング
    for {
		//接続要求を処理する
		conn, err := listener.Accept()
		defer conn.Close()

		log.Print("New Socket Connection Established From ",conn.RemoteAddr())
        if err != nil {
			connFinished(conn,err)
            continue
		}

		//接続が確率したら別スレッドで処理を開始する
        go handler(conn)

    }
}

func connFinished(conn net.Conn,err error){
	log.Print("Socket Connection ",conn.RemoteAddr()," Finished by Error ",err.Error())
}

func handler(conn net.Conn){
	//一番最初の一行を名前として扱う
	r := bufio.NewReader(conn)
	name,err := r.ReadString('\n')
	if err!=nil{
		connFinished(conn,err)
		return
	}
	//CSV的に処理するのでカンマを消す
	name = strings.Replace(name," ","",-1)
	//最後に改行があるので削除する
	name = name[:len(name)-1]
	//デバイスのインスタンスを作成
	d:=device.NewDevice(name,conn)
	//管理対象に追加
	Devices[name]=d
	log.Print("New Device Defined : ",name,"(",conn.RemoteAddr(),")")
	for {
		//入力された文字をdeviceで処理する
		str,err := r.ReadString('\n')
		if err!=nil{
			connFinished(conn,err)
			return
		}
		d.Parse(str[:len(str)-1])
	}
}