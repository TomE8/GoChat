package common

type ReadRoomManager struct {
	RoomName string
	RespChan chan *Room
}

type DelRoomManager struct {
	RoomName string
}

type RoomManger struct {
	Read chan *ReadRoomManager
	Del  chan *DelRoomManager
}

func (m ReadRoomManager) GetRoomName() string {
	return m.RoomName
}

func (m DelRoomManager) GetRoomName() string {
	return m.RoomName
}

func (r RoomManger) GetRoom(roomName string) *Room {
	respCha := make(chan *Room)
	message := ReadRoomManager{RoomName: roomName, RespChan: respCha}
	r.Read <- &message
	return <-respCha
}
