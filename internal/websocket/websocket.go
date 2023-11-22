package websocket

import (
    "encoding/json"
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

type Message struct {
    Type     string `json:"type"`
    Username string `json:"username"`
    Text string `json:"text"`
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
        message, err := wsClient.Conn.ReadMessage()
        if err != nil {
            if err == io.EOF {
                fmt.Println("Client disconnected")
            } else {
                fmt.Println("Error reading message:", err)
            }
            break
        }

        var wsResponse Message
        err = json.Unmarshal([]byte(message), &wsResponse)
        if err != nil {
            fmt.Println("Error unmarshalling client message", err)

            wsResponse = Message{
                Text: string(message),
                Type: "message",
            }
        }

        msg := chat.Message{
            Username:  wsClient.Username,
            Text:      wsResponse.Text,
            Type:      wsResponse.Type,
            Timestamp: time.Now(),
        }

        chatServer.Broadcast(msg)
    }
}

func (wc *WebSocketClient) Send(msg chat.Message) error{
    jsonData, err := json.Marshal(msg)
    if err != nil {
        fmt.Println("Error marshalling message to JSON:", err)
        return err
    }

    err = wc.Conn.WriteMessage(jsonData)
    if err != nil {
        fmt.Println("Error sending message:", err)
        return err
    }

    return nil
}
