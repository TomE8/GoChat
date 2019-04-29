package common

const (
	ReadMode        = 0
	ReadAllMode     = 1
	ReadFromEndMode = 2

	AppendMessage = 3
	EnterRoom     = 4
	ExitRoom      = 99

	Port = "8200"
	Ip = "127.0.0.1"
)

type ClientCommand struct {
	RoomName   string
	Flag       int
	NumOfLines int
	Content    string
}

type ServerReturn struct {
	Data []string
}

func (s ServerReturn) GetData() []string{
	return s.Data
}

func (c ClientCommand) GetRoomName() string{
	return c.RoomName
}

func (c ClientCommand) GetFlag() int {
	return c.Flag
}

func (c ClientCommand) GetNumOfLines() int {
	return c.NumOfLines
}

func (c ClientCommand) GetContent() string {
	return c.Content
}
