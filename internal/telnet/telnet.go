package telnet

import (
	"bufio"
	"fmt"
    "log"
	"net"
	"time"

	"os"

	"github.com/adriankuklinski/go-chat/internal/chat"
	"github.com/adriankuklinski/go-chat/pkg/utils"
)

type TelnetClient struct {
	ID       string
	Username string
	Conn     net.Conn
}

func (tnClient *TelnetClient) Send(msg chat.Message) error {
	formattedMessage := fmt.Sprintf("%s [%s]: %s\n", msg.Timestamp.Format("15:04:05"), msg.Username, msg.Text)

	_, err := tnClient.Conn.Write([]byte(formattedMessage))
	if err != nil {
		fmt.Println("Error sending message to Telnet client:", err)
		tnClient.Conn.Close()
        return err
	}

    return nil
}

func StartTelnetServer(addr string, chatServer *chat.Server) {
    fmt.Println("Starting Telnet server on", addr)

	ln, err := net.Listen("tcp", addr)
	if err != nil {
        fmt.Println("Error starting telnet:", err)
	}

	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		tnClient := &TelnetClient{
			ID:       utils.GenerateUniqueID(),
			Username: "default",
			Conn:     conn,
		}

		chatServer.AddClient(tnClient)
		go handleConnection(tnClient, chatServer)
        log.Println("New connection from:", conn.RemoteAddr())
	}
}

func handleConnection(tnClient *TelnetClient, chatServer *chat.Server) {
	defer tnClient.Conn.Close()
	scanner := bufio.NewScanner(tnClient.Conn)
	for scanner.Scan() {
		text := scanner.Text()
		msg := chat.Message{
			Username:  tnClient.Username,
			Text:      text,
            Type:      "message",
			Timestamp: time.Now(),
		}

		chatServer.Broadcast(msg)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stdout, []any{"Telnet scanner error: %v", err}...)
	}

	chatServer.RemoveClient(tnClient)
}
