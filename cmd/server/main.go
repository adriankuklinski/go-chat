package main

import (
    "embed"
    "net/http"

    "github.com/adriankuklinski/go-chat/internal/chat"
    "github.com/adriankuklinski/go-chat/internal/telnet"
    "github.com/adriankuklinski/go-chat/internal/websocket"
)

var content embed.FS

func main() {
    chatServer := chat.NewServer()

    http.Handle("/", http.FileServer(http.FS(content)))

    go telnet.StartTelnetServer(":8081", chatServer)
    go websocket.StartWebSocketServer(":8080", chatServer)

    select {}
}
