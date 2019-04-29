package main

import (
	"Chat/Common"
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

func PrintManyLines(data []string){
	if len(data)==0{
		return
	}
	stringSlices := strings.Join(data, "\n")
	fmt.Println(stringSlices)
}

func ReceiveFromServer(conn net.Conn) []string{
	message, err := bufio.NewReader(conn).ReadBytes('}')
	if err != nil {
		log.Fatal(err)
	}
	var parsedMessage Common.ServerReturn
	err = json.Unmarshal(message, &parsedMessage)
	if err != nil {
		log.Fatal(err)
	}
	return parsedMessage.GetData()
}

func SendToServer(conn net.Conn, serverMessage Common.ClientCommand) {
	messageToSend, err := json.Marshal(serverMessage)
	if err != nil {
		log.Fatal(err)
	}
	if _, err = conn.Write(messageToSend); err != nil {
		log.Fatal(err)
	}
}

func ConnectionHandler(writeChan chan string, exitChan chan bool) {
	var serverMessage Common.ClientCommand
	var NumOfMessages = 0
	conn, err := net.Dial("tcp", Common.Ip+":"+Common.Port)
	if err != nil {
		panic(err)
	}
	text := fmt.Sprintf("---> %s has enterd the room", UserName)
	serverMessage = Common.ClientCommand{RoomName: RoomName, Flag: Common.EnterRoom, NumOfLines: 0, Content: text}
	SendToServer(conn, serverMessage)
	for {
		select {
		case message := <-writeChan:
			serverMessage = Common.ClientCommand{RoomName: RoomName, Flag: Common.AppendMessage, NumOfLines: 0, Content: message}
			SendToServer(conn, serverMessage)
			NumOfMessages++
		case <-time.After(RefreshTimeout):
			serverMessage = Common.ClientCommand{RoomName: RoomName, Flag: Common.ReadFromEndMode, NumOfLines: NumOfMessages, Content: ""}
			SendToServer(conn, serverMessage)
			data := ReceiveFromServer(conn)
			PrintManyLines(data)
			NumOfMessages = NumOfMessages + len(data)
		case <-exitChan:
			text := fmt.Sprintf("---> %s has left the room", UserName)
			serverMessage = Common.ClientCommand{RoomName: RoomName, Flag: Common.ExitRoom, NumOfLines: 0, Content: text}
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
