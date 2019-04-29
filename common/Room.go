package common

type ReadRoomMessage struct {
	ReadMode   int
	NumOfLines int
	RespChan   chan []string
}

type WriteRoomMessage struct {
	Flag    int
	Message string
}

type CtrlRoomMessage struct {
	Flag int
}

type Room struct {
	RoomName     string
	Read         chan *ReadRoomMessage
	Write        chan *WriteRoomMessage
	Ctrl         chan *CtrlRoomMessage
	Del          chan *DelRoomManager
	NumOfClients int
}

func (r *Room) GetRoomName() string{
	return r.RoomName
}

func (r *Room) GetReadChan() chan *ReadRoomMessage {
	return r.Read
}

func (r *Room) GetNumOfClients() int {
	return r.NumOfClients
}

func (r *Room) AddClient() {
	r.NumOfClients++
}

func (r *Room) RemoveClient() {
	r.NumOfClients--
}
