package client

import (
	"Chat/common"
	"bufio"
	"encoding/json"
	"log"
	"net"
)

func SendDataToClient(c net.Conn, data []string) {
	message := common.ServerResponse{Data: data}
	messageToSend, err := json.Marshal(message)
	if err != nil {
		log.Fatal(err)
	}
	if _, err = c.Write(messageToSend); err != nil {
		log.Fatal(err)
	}
}

func ParseClientCommand(c net.Conn) common.ClientCommand {

	message, err := bufio.NewReader(c).ReadBytes('}')
	if err != nil {
		log.Fatal(err)
	}
	var parsedMessage common.ClientCommand
	err = json.Unmarshal(message, &parsedMessage)
	if err != nil {
		log.Fatal(err)
	}
	return parsedMessage
}

func Handler(c net.Conn, roomManger *common.RoomManger) {
	for {
		clientCommand := ParseClientCommand(c)
		room := roomManger.GetRoom(clientCommand.GetRoomName())
		switch clientCommand.GetFlag() {
		case common.ExitRoom:
			room.Write <- &common.WriteRoomMessage{Flag: common.AppendMessage, Message: clientCommand.GetContent()}
			room.Ctrl <- &common.CtrlRoomMessage{Flag: common.ExitRoom}
			return
		case common.EnterRoom:
			room.Write <- &common.WriteRoomMessage{Flag: common.AppendMessage, Message: clientCommand.GetContent()}
			room.Ctrl <- &common.CtrlRoomMessage{Flag: common.EnterRoom}
		case common.AppendMessage:
			room.Write <- &common.WriteRoomMessage{Flag: common.AppendMessage, Message: clientCommand.GetContent()}
		case common.ReadAllMode, common.ReadFromEndMode:
			RespChan := make(chan []string)
			room.Read <- &common.ReadRoomMessage{ReadMode: common.ReadMode, NumOfLines: clientCommand.GetNumOfLines(), RespChan: RespChan}
			messageArray := <-RespChan
			if clientCommand.GetFlag() == common.ReadFromEndMode {
				messageArray = messageArray[clientCommand.GetNumOfLines():]
			}
			SendDataToClient(c, messageArray)
			close(RespChan)
		}
	}
}
