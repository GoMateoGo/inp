package gateway

import (
	"net"

	"github.com/GoMateoGo/inp/pkg/logger"
)

type Geteway struct {
	ListenAddr string // 监听地址
	sessionMgr *SessionManager
}

func NewGeteway(sMgr *SessionManager) *Geteway {
	return &Geteway{
		sessionMgr: sMgr,
	}
}

func (g *Geteway) ListenAndServer() error {
	listener, err := net.Listen("tcp", g.ListenAddr)
	if err != nil {
		return err
	}

	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}

		go g.handleTCPConnection(conn)
	}
}

func (g *Geteway) handleTCPConnection(conn net.Conn) {
	defer conn.Close()

	handShakeReq := &HandShakeReq{}
	if err := handShakeReq.Dcode(conn); err != nil {
		logger.Error("decode handshake fail:%v", err)
		return
	}

	// 创建session
}
