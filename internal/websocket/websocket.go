package websocket

import (
    "fmt"
    "io"
    "log"
    "net"
    "net/http"
    "os"
    "time"

    "github.com/adriankuklinski/go-chat/internal/chat"
    "github.com/adriankuklinski/go-chat/pkg/utils"
    ws "github.com/adriankuklinski/go-chat/pkg/websocket"
)

type WebSocketClient struct {
    ID       string
    Username string
    Conn     *ws.WebSocketConnection
}

func StartWebSocketServer(addr string, chatServer *chat.Server) {
    fmt.Println("WebSocket server listening on", addr)

    http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
        conn, err := ws.UpgradeToWebSocket(w, r)
        if err != nil {
            fmt.Println("WebSocket upgrade failed:", err)
            return
        }

        handleConnection(conn, chatServer)
    })

    err := http.ListenAndServe(addr, nil)
    if err != nil {
		fmt.Fprintln(os.Stdout, []any{"Websocket ListenAndServe: %v", err}...)
    }
}

func handleConnection(conn net.Conn, chatServer *chat.Server) {
    defer conn.Close()

    wsConn := &ws.WebSocketConnection{Conn: conn}
    wsClient := &WebSocketClient{
        ID:       utils.GenerateUniqueID(),
        Username: "default",
        Conn:     wsConn,
    }

    chatServer.AddClient(wsClient)

    log.Println("New connection from:", conn.RemoteAddr())

    for {
        fmt.Println("Attempting to read message")
        message, err := wsClient.Conn.ReadMessage()
        if err != nil {
            if err == io.EOF {
                fmt.Println("Client disconnected")
            } else {
                fmt.Println("Error reading message:", err)
            }
            break
        }

        msg := chat.Message{
            Username:  wsClient.Username,
            Text:      string(message),
            Timestamp: time.Now(),
        }

        chatServer.Broadcast(msg)
    }
}

func (wc *WebSocketClient) Send(msg chat.Message) {
    formattedMessage := fmt.Sprintf("%s [%s]: %s\n", msg.Timestamp.Format("15:04:05"), msg.Username, msg.Text)

    err := wc.Conn.WriteMessage([]byte(formattedMessage))
    if err != nil {
        fmt.Println("Error sending message:", err)
    }
}
