package gateway

type ClientInfo struct {
	ClientID         string // 客户端id
	PublicProtocol   string // 公网协议
	PublicIP         string // 公网ip
	PublicPort       uint16 // 公网端口
	InternalProtocol string // 内网协议
	InternalIP       string // 内网ip
	InternalPort     uint16 // 内网端口
}