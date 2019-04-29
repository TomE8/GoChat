package room

import (
	"Chat/common"
	"fmt"
	"time"
)

const RefreshTimeout  = 5 * time.Second // TODO: move this out of here

func NotifyCloseStatusRoomManager(room *common.Room){
	room.Del<-&common.DelRoomManager{RoomName: room.GetRoomName()}
}

func CloseRoomChannels(room *common.Room){
	close(room.Read)
	close(room.Write)
	close(room.Ctrl)
}

func CtrlHandle(room *common.Room, ctrlChan *common.CtrlRoomMessage) {
	switch ctrlChan.Flag {
	case common.EnterRoom:
		room.AddClient()
	case common.ExitRoom:
		room.RemoveClient()
	}
}

func Handler(room *common.Room) {
	var messageArray []string
	for {
		select {
		case readChan := <-room.Read:
			readChan.RespChan <- messageArray
		case writeChan := <-room.Write:
			messageArray = append(messageArray, writeChan.Message)
		case ctrlChan := <-room.Ctrl:
			CtrlHandle(room, ctrlChan)
		case <-time.After(RefreshTimeout):
			if room.GetNumOfClients()==0 {
				fmt.Printf("room %s is closed\n",room.GetRoomName())
				time.Sleep(time.Second)
				CloseRoomChannels(room)
				NotifyCloseStatusRoomManager(room)
				return
			}
		}
	}
}
