package common

type ReadRoomChan struct {
	ReadMode   int
	NumOfLines int
	RespChan   chan []string
}

type WriteRoomChan struct {
	Flag    int
	Message string
}

type CtrlRoomChan struct {
	Flag int
}

type Room struct {
	RoomName     string
	Read         chan *ReadRoomChan
	Write        chan *WriteRoomChan
	Ctrl         chan *CtrlRoomChan
	Del          chan *DelRoomManager
	NumOfClients int
}

func (r *Room) GetRoomName() string{
	return r.RoomName
}

func (r *Room) GetReadChan() chan *ReadRoomChan{
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
