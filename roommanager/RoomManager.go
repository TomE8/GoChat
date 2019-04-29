package roommanager

import (
	"Chat/room"
	"Chat/common"
)

func createNewRoom(c chan *common.DelRoomManager, roomName string) *common.Room {
	roomReadChan := make(chan *common.ReadRoomMessage)
	roomWriteChan := make(chan *common.WriteRoomMessage)
	roomCtrlChan := make(chan *common.CtrlRoomMessage)
	return &common.Room{RoomName: roomName,Read: roomReadChan, Write: roomWriteChan, Ctrl: roomCtrlChan, Del: c, NumOfClients: 0}
}

func appendNewRoom(channelPerRoomMapping map[string]*common.Room, c chan *common.DelRoomManager, roomName string) *common.Room {
	pNewRoom := createNewRoom(c,roomName)
	channelPerRoomMapping[roomName] = pNewRoom
	return pNewRoom
}

func Handler(roomManger common.RoomManger) {
	channelPerRoomMapping := make(map[string]*common.Room)
	for {
		select {
		case readChan := <-roomManger.Read:
			pRoom, ok := channelPerRoomMapping[readChan.GetRoomName()]
			if !ok {
				pRoom = appendNewRoom(channelPerRoomMapping, roomManger.Del, readChan.GetRoomName())
				go room.Handler(pRoom)
			}
			readChan.RespChan <- pRoom
		case delChan := <-roomManger.Del:
			_, ok := channelPerRoomMapping[delChan.GetRoomName()]
			if !ok {
				return
			}
			delete(channelPerRoomMapping, delChan.GetRoomName())
		}
	}
}
