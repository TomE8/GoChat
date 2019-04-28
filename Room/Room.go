package Room

import (
	"Chat/Common"
	"fmt"
	"time"
)

const RefreshTimeout  = 5 * time.Second // TODO: move this out of here

func NotifyCloseStatusRoomManager(room *Common.Room){
	room.Del<-&Common.DelRoomManager{RoomName: room.GetRoomName()}
}

func CloseRoomChannels(room *Common.Room){
	close(room.Read)
	close(room.Write)
	close(room.Ctrl)
}

func CtrlHandle(room *Common.Room, ctrlChan *Common.CtrlRoomChan) {
	switch ctrlChan.Flag {
	case Common.EnterRoom:
		room.AddClient()
	case Common.ExitRoom:
		room.RemoveClient()
	}
}

func Handler(room *Common.Room) {
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
