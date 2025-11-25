package gateway

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
)

const (
	cmdPP = 0x0 // proxy protocol
	cmdHS = 0x1 // handshake
)

// 私有协议

type ProxyProtocol struct {
	CLientID         string // 客户端id
	PublicProtocol   string // 公共协议
	PublicIP         string // 公共ip
	PublicPort       uint16 // 公共端口
	InternalProtocol string // 内网协议
	InternalIP       string // 内网ip
	InternalPort     uint16 // 内网端口
}

// 1byte version
// 1byte cmd
// 2bytes length
// length body
func (pp *ProxyProtocol) Encode() ([]byte, error) {
	hdr := make([]byte, 4)
	hdr[0] = cmdPP // version
	hdr[1] = 0x0   // cmd

	body, err := json.Marshal(pp)
	if err != nil {
		return nil, err
	}

	binary.BigEndian.PutUint16(hdr[2:4], uint16(len(body)))
	return append(hdr, body...), nil
}

type HandShakeReq struct {
	ClientID string // 客户端id
}

func (req *HandShakeReq) Encode() ([]byte, error) {
	hdr := make([]byte, 4)
	hdr[0] = 0x0   // version
	hdr[1] = cmdHS // cmd

	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	binary.BigEndian.PutUint16(hdr[2:4], uint16(len(body)))
	return append(hdr, body...), nil
}

func (req *HandShakeReq) Dcode(reader io.Reader) error {
	hdr := make([]byte, 4)
	len, err := io.ReadFull(reader, hdr)
	if err != nil {
		return err
	}
	if len < 4 {
		return fmt.Errorf("req pacakge len <4")
	}

	cmd := hdr[1]
	if cmd != cmdHS {
		return fmt.Errorf("invalid handshake cmd.")
	}

	bodyLen := binary.BigEndian.Uint16(hdr[2:4])

	body := make([]byte, bodyLen)
	_, err = io.ReadFull(reader, body)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(body, req); err != nil {
		return err
	}
	
	return nil
}
