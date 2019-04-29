package main

import (
	"Chat/client"
	"Chat/common"
	"Chat/roommanager"
	"log"
	"net"
)

func main() {

	RMReadChan := make(chan *common.ReadRoomManager)
	RMDelChan := make(chan *common.DelRoomManager)
	roomManger := common.RoomManger{Read: RMReadChan, Del: RMDelChan}

	go roommanager.Handler(roomManger)

	l, err := net.Listen("tcp", ":"+common.Port)
	if err != nil {
		log.Fatal(err)
	}

	for {
		c, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go client.Handler(c, &roomManger)
	}
	if err := l.Close(); err != nil {
		log.Fatal(err)
	}

}
