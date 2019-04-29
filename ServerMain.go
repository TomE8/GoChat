package main

import (
	"Chat/Client"
	"Chat/Common"
	"Chat/RoomManager"
	"log"
	"net"
)

func main() {

	RMReadChan := make(chan *Common.ReadRoomManager)
	RMDelChan := make(chan *Common.DelRoomManager)
	roomManger := Common.RoomManger{Read: RMReadChan, Del: RMDelChan}

	go RoomManager.Handler(roomManger)

	l, err := net.Listen("tcp", ":"+Common.Port)
	if err != nil {
		log.Fatal(err)
	}

	for {
		c, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go Client.Handler(c, &roomManger)
	}
	if err := l.Close(); err != nil {
		log.Fatal(err)
	}

}
