package sse

import (
	"github.com/doorman2137/betonz-go/internal/db"
	"github.com/jackc/pgx/v5/pgtype"
)

type Connection struct {
	User           db.User
	MessageChannel chan string
}

type EventServer struct {
	connections map[[16]byte]Connection
}

func NewServer() *EventServer {
	return &EventServer{connections: make(map[[16]byte]Connection)}
}

func (s *EventServer) Subscribe(id [16]byte, user db.User) Connection {
	connection := Connection{
		User:           user,
		MessageChannel: make(chan string),
	}
	s.connections[id] = connection
	return connection
}

func (s *EventServer) Unsubscribe(id [16]byte) {
	delete(s.connections, id)
}

func (s *EventServer) Notify(userId pgtype.UUID, message string) {
	for _, connection := range s.connections {
		if connection.User.ID.Bytes == userId.Bytes {
			connection.MessageChannel <- message
		}
	}
}

func (s *EventServer) NotifyAdmins(message string) {
	for _, connection := range s.connections {
		if connection.User.Role == db.RoleADMIN || connection.User.Role == db.RoleSUPERADMIN {
			connection.MessageChannel <- message
		}
	}
}
