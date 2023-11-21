package chat

import (
    "fmt"
    "strings"
    "time"
)

type Client interface {
    Send(Message) error
}

type Message struct {
    Username  string
    Text      string
    Timestamp time.Time
    Type      string
}

type Server struct {
    clients []Client
    typingClients map[string]bool
}

func NewServer() *Server {
    return &Server{
        clients: make([]Client, 0),
        typingClients: make(map[string]bool),
    }
}

func (s *Server) Broadcast(message Message) {
    for _, client := range s.clients {
        client.Send(message)
    }
}

func (s *Server) AddClient(client Client) {
    s.clients = append(s.clients, client)
}

func (s *Server) RemoveClient(client Client) {
    for i, client := range s.clients {
        if client == client {
            s.clients = append(s.clients[:i], s.clients[i+1:]...)
            break
        }
    }
}

func (s *Server) UpdateTypingStatus(clientID string, isTyping bool) {
    s.typingClients[clientID] = isTyping
    s.sendTypingIndicator()
}

func (s *Server) sendTypingIndicator() {
    var typingUsers []string
    for clientID, isTyping := range s.typingClients {
        if isTyping {
            typingUsers = append(typingUsers, clientID)
        }
    }

    var typingMessage string
    if len(typingUsers) > 1 {
        typingMessage = strings.Join(typingUsers, ", ") + " are typing..."
    } else if len(typingUsers) == 1 {
        typingMessage = typingUsers[0] + " is typing..."
    }

    if typingMessage != "" {
        message := Message{
            Username: "System",
            Text:     typingMessage,
            Timestamp: time.Now(),
            Type: "typing",
        }

        for _, client := range s.clients {
            if err := client.Send(message); err != nil {
                fmt.Println("Error sending typing indicator:", err)
            }
        }
    }
}
