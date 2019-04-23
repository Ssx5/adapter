package device

import (
	"fmt"

	"adapter/device/server"
)

func NewDeviceServer(d *DeviceInfo) (ds server.DeviceServer, err error) {
	switch d.Protocol {
	case "tcp-server":
		ds = NewTcpServer(d)
	case "iwsn":
		ds = NewUsbServer(d)
	default:
		err = fmt.Errorf("Unsupport protocol %s", d.Protocol)
	}
	return
}

func NewTcpServer(d *DeviceInfo) *server.TcpServer {
	return &server.TcpServer{
		Type:     d.Type,
		Protocol: d.Protocol,
		Name:     d.Name,
		Address:  d.Parameter["address"],
		Timeout:  d.Parameter["timeout"],
		State:    server.StateOffline,
	}
}

func NewUsbServer(d *DeviceInfo) *server.UsbServer {
	return &server.UsbServer{
		Type:      d.Type,
		Protocol:  d.Protocol,
		Name:      d.Name,
		Interface: d.Parameter["interface"],
		EthFilter: d.Parameter["ethfilter"],
		Timeout:   d.Parameter["timeout"],
		State:     server.StateOffline,
	}
}
