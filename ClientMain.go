package main

import (
	"Chat/common"
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

const (
	ExitWord       = "exit()"
	UserName       = "TOM"
	RoomName       = "temp"
	RefreshTimeout = 1 * time.Second
)

func PrintManyLines(data []string) {
	if len(data) == 0 {
		return
	}
	stringSlices := strings.Join(data, "\n")
	fmt.Println(stringSlices)
}

func ReceiveFromServer(conn net.Conn) []string {
	message, err := bufio.NewReader(conn).ReadBytes('}')
	if err != nil {
		log.Fatal(err)
	}
	var parsedMessage common.ServerResponse
	err = json.Unmarshal(message, &parsedMessage)
	if err != nil {
		log.Fatal(err)
	}
	return parsedMessage.GetData()
}

func SendToServer(conn net.Conn, serverMessage common.ClientCommand) {
	messageToSend, err := json.Marshal(serverMessage)
	if err != nil {
		log.Fatal(err)
	}
	if _, err = conn.Write(messageToSend); err != nil {
		log.Fatal(err)
	}
}

func ConnectionHandler(writeChan chan string, exitChan chan bool) {
	var serverMessage common.ClientCommand
	var NumOfMessages = 0
	conn, err := net.Dial("tcp", common.Ip+":"+common.Port)
	if err != nil {
		panic(err)
	}
	text := fmt.Sprintf("---> %s has enterd the room", UserName)
	serverMessage = common.ClientCommand{RoomName: RoomName, Flag: common.EnterRoom, NumOfLines: 0, Content: text}
	SendToServer(conn, serverMessage)
	for {
		select {
		case message := <-writeChan:
			serverMessage = common.ClientCommand{RoomName: RoomName, Flag: common.AppendMessage, NumOfLines: 0, Content: message}
			SendToServer(conn, serverMessage)
			NumOfMessages++
		case <-time.After(RefreshTimeout):
			serverMessage = common.ClientCommand{RoomName: RoomName, Flag: common.ReadFromEndMode, NumOfLines: NumOfMessages, Content: ""}
			SendToServer(conn, serverMessage)
			data := ReceiveFromServer(conn)
			PrintManyLines(data)
			NumOfMessages = NumOfMessages + len(data)
		case <-exitChan:
			text := fmt.Sprintf("---> %s has left the room", UserName)
			serverMessage = common.ClientCommand{RoomName: RoomName, Flag: common.ExitRoom, NumOfLines: 0, Content: text}
			SendToServer(conn, serverMessage)
			exitChan <- true
			return
		}
	}
}

func main() {
	exitChan := make(chan bool)
	writeChan := make(chan string)

	go ConnectionHandler(writeChan, exitChan)

	reader := bufio.NewReader(os.Stdin)
	for {
		text, _ := reader.ReadString('\n')
		text = text[:len(text)-1] // remove \n

		if text == ExitWord {
			exitChan <- true
			time.Sleep(time.Second)
			<-exitChan
			return
		}
		text = fmt.Sprintf("%s: %s", UserName, text)
		writeChan <- text

	}
}
