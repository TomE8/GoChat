package Client

import (
	"Chat/Common"
	"bufio"
	"encoding/json"
	"log"
	"net"
)

type RoomI interface {
	GetRoomName() string
	GetNumOfClients() int
	AddClient()
	RemoveClient()
}

type RoomManagerI interface {
	GetRoom(roomName string) *RoomI
}

func SendDataToClient(c net.Conn, data []string) {
	message := Common.ServerReturn{Data: data}
	messageToSend, err := json.Marshal(message)
	if err != nil {
		log.Fatal(err)
	}
	if _, err = c.Write(messageToSend); err != nil {
		log.Fatal(err)
	}
}

func ParseClientCommand(c net.Conn) Common.ClientCommand {

	message, err := bufio.NewReader(c).ReadBytes('}')
	if err != nil {
		panic(err)
	}
	var parsedMessage Common.ClientCommand
	err = json.Unmarshal(message, &parsedMessage)
	if err != nil {
		panic(err)
	}
	return parsedMessage
}

func Handler(c net.Conn, roomManger *Common.RoomManger) {
	for {
		clientCommand := ParseClientCommand(c)
		room := roomManger.GetRoom(clientCommand.GetRoomName())
		switch clientCommand.GetFlag() {
		case Common.ExitRoom:
			room.Write <- &Common.WriteRoomChan{Flag: Common.AppendMessage, Message: clientCommand.GetContent()}
			room.Ctrl <- &Common.CtrlRoomChan{Flag: Common.ExitRoom}
			return
		case Common.EnterRoom:
			room.Write <- &Common.WriteRoomChan{Flag: Common.AppendMessage, Message: clientCommand.GetContent()}
			room.Ctrl <- &Common.CtrlRoomChan{Flag: Common.EnterRoom}
		case Common.AppendMessage:
			room.Write <- &Common.WriteRoomChan{Flag: Common.AppendMessage, Message: clientCommand.GetContent()}
		case Common.ReadAllMode, Common.ReadFromEndMode:
			RespChan := make(chan []string)
			room.Read <- &Common.ReadRoomChan{ReadMode: Common.ReadMode, NumOfLines: clientCommand.GetNumOfLines(), RespChan: RespChan}
			messageArray := <-RespChan
			if clientCommand.GetFlag() == Common.ReadFromEndMode {
				messageArray = messageArray[clientCommand.GetNumOfLines():]
			}
			SendDataToClient(c, messageArray)
			close(RespChan)
		}
	}
}
