package gateway

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"github.com/GoMateoGo/inp/pkg/logger"
)

var (
	writeTimeout = time.Second * 3
)

type Listener struct {
	pp          *ProxyProtocol  // 私有协议
	sessionMgr  *SessionManager // session管理器
	closeOnce   sync.Once       // 关闭一次
	close       chan struct{}   // 关闭通道
	tcpListener net.Listener    // tcp监听器
}

func NewListener(pp *ProxyProtocol, sessionMgr *SessionManager) *Listener {
	return &Listener{
		close:      make(chan struct{}),
		sessionMgr: sessionMgr,
		pp:         pp,
	}
}

func (l *Listener) ListenAndServer() error {
	switch l.pp.PublicProtocol {
	case "tcp":
		return l.listenAndServerTCP()
	default:
		return fmt.Errorf("TODO:..")
	}
}

func (l *Listener) Close() {
	l.closeOnce.Do(func() {
		close(l.close)
		if l.tcpListener != nil {
			l.tcpListener.Close()
		}
	})
}

func (l *Listener) listenAndServerTCP() error {
	listenAddr := fmt.Sprintf("%s:%d", l.pp.PublicIP, l.pp.PublicPort)
	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return err
	}

	defer listener.Close()
	l.tcpListener = listener

	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}

		go l.handleTCPConnection(conn)
	}
}

func (l *Listener) handleTCPConnection(conn net.Conn) {
	defer conn.Close()
	// 查询session
	tunnelConn, err := l.sessionMgr.GetSessionByClientID(l.pp.CLientID)
	if err != nil {
		logger.Error("get session for client %s failed", l.pp.CLientID)
		return
	}
	defer tunnelConn.Close()

	// 封装proxyProtocol
	ppBody, err := l.pp.Encode()
	if err != nil {
		logger.Error("encode proxy protocol failed: %v", err)
		return
	}

	tunnelConn.SetWriteDeadline(time.Now().Add(writeTimeout))

	_, err = tunnelConn.Write(ppBody)
	tunnelConn.SetWriteDeadline(time.Time{})
	if err != nil {
		logger.Error("send proxy protocol to client %s failed: %v", l.pp.CLientID, err)
		return
	}

	// 双向数据拷贝
	io.Copy(conn, tunnelConn)
	go io.Copy(tunnelConn, conn)
}
