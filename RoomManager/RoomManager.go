package RoomManager

import (
	"Chat/Room"
	"Chat/Common"
)

func createNewRoom(c chan *Common.DelRoomManager, roomName string) *Common.Room {
	roomReadChan := make(chan *Common.ReadRoomChan)
	roomWriteChan := make(chan *Common.WriteRoomChan)
	roomCtrlChan := make(chan *Common.CtrlRoomChan)
	return &Common.Room{RoomName: roomName,Read: roomReadChan, Write: roomWriteChan, Ctrl: roomCtrlChan, Del: c, NumOfClients: 0}
}

func appendNewRoom(channelPerRoomMapping map[string]*Common.Room, c chan *Common.DelRoomManager, roomName string) *Common.Room {
	pNewRoom := createNewRoom(c,roomName)
	channelPerRoomMapping[roomName] = pNewRoom
	return pNewRoom
}

func Handler(roomManger Common.RoomManger) {
	channelPerRoomMapping := make(map[string]*Common.Room)
	for {
		select {
		case readChan := <-roomManger.Read:
			pRoom, ok := channelPerRoomMapping[readChan.GetRoomName()]
			if !ok {
				pRoom = appendNewRoom(channelPerRoomMapping, roomManger.Del, readChan.GetRoomName())
				go Room.Handler(pRoom)
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
