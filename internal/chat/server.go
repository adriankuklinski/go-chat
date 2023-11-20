package chat

import (
    "time"
)

type Client interface {
    Send(Message)
}

type Message struct {
    Username  string
    Text      string
    Timestamp time.Time
}

type Server struct {
    clients []Client
}

func NewServer() *Server {
    return &Server{
        clients: make([]Client, 0),
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
