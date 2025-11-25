package gateway

import (
	"fmt"
	"net"
	"sync"

	"github.com/xtaci/smux"
)

type Session struct {
	ClientID   string
	Connection smux.Session
}

type SessionManager struct {
	session  sync.Mutex
	sessions map[string]*Session
}

func NewSessionManager() *SessionManager {
	return &SessionManager{
		sessions: make(map[string]*Session),
	}
}

func (s *SessionManager) GetSessionByClientID(clientID string) (net.Conn, error) {
	s.session.Lock()
	defer s.session.Unlock()
	session := s.sessions[clientID]
	if session == nil {
		return nil, fmt.Errorf("client %s not connected", clientID)
	}

	stream, err := session.Connection.OpenStream()
	if err != nil {
		return nil, fmt.Errorf("open stream for client %s failed: %v", clientID, err)
	}
	return stream, nil
}
