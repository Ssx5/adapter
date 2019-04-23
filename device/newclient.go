package device

import (
	"fmt"
	"strconv"

	"adapter/device/client"
)

func NewDeviceClient(d *DeviceInfo) (dc client.DeviceClient, err error) {
	switch d.Protocol {
	case "virtual":
		dc, err = NewSomeClient(d)
	case "tcp-client":
		dc, err = NewTcpClient(d)
	case "modbus-tcp":
		dc, err = NewModbusTcpClient(d)
	case "modus-rtu":
		dc, err = NewModbusRtuClient(d)
	case "serial":
		dc, err = NewSerialClient(d)
	default:
		err = fmt.Errorf("Unsupport protocol %s", d.Protocol)
	}
	return
}

func NewModbusRtuClient(d *DeviceInfo) (*client.ModbusRtuClient, error) {
	var err error
	if _, ok := d.Parameter["baudrate"]; !ok {
		return nil, fmt.Errorf("section 'parameter' should have field 'baudrate'")
	}
	baudrate, err := strconv.Atoi(d.Parameter["baudrate"])
	if err != nil {
		return nil, err
	}
	if _, ok := d.Parameter["databits"]; !ok {
		return nil, fmt.Errorf("section 'parameter' should have field 'databits'")
	}
	databits, err := strconv.Atoi(d.Parameter["databits"])
	if err != nil {
		return nil, err
	}
	if _, ok := d.Parameter["stopbits"]; !ok {
		return nil, fmt.Errorf("section 'parameter' should have field 'stopbits'")
	}
	stopbits, err := strconv.Atoi(d.Parameter["stopbits"])
	if err != nil {
		return nil, err
	}
	if _, ok := d.Parameter["slaveid"]; !ok {
		return nil, fmt.Errorf("section 'parameter' should have field 'slaveid'")
	}
	slaveid, err := strconv.Atoi(d.Parameter["slaveid"])
	if err != nil {
		return nil, err
	}
	if _, ok := d.Parameter["path"]; !ok {
		return nil, fmt.Errorf("section 'parameter' should have field 'path'")
	}
	if _, ok := d.Parameter["parity"]; !ok {
		return nil, fmt.Errorf("section 'parameter' should have field 'parity'")
	}
	for _, c := range d.Commands {
		err = c.Check([]string{"funCode", "startAddr", "quantity"})
		if err != nil {
			return nil, err
		}
	}
	return &client.ModbusRtuClient{
		Type:     d.Type,
		Protocol: d.Protocol,
		Name:     d.Name,
		Path:     d.Parameter["path"],
		BaudRate: baudrate,
		DataBits: databits,
		Parity:   d.Parameter["parity"],
		StopBits: stopbits,
		SlaveId:  byte(slaveid),
		State:    client.StateOffline,
		Commands: d.Commands,
	}, nil
}

func NewModbusTcpClient(d *DeviceInfo) (*client.ModbusTcpClient, error) {
	if _, ok := d.Parameter["slaveid"]; !ok {
		return nil, fmt.Errorf("section 'parameter' should have field 'slaveid'")
	}
	slaveid, err := strconv.Atoi(d.Parameter["slaveid"])
	if err != nil {
		return nil, err
	}
	if _, ok := d.Parameter["address"]; !ok {
		return nil, fmt.Errorf("section 'parameter' should have field 'address'")
	}
	for _, c := range d.Commands {
		err = c.Check([]string{"funCode", "startAddr", "quantity"})
		if err != nil {
			return nil, err
		}
	}
	return &client.ModbusTcpClient{
		Type:     d.Type,
		Protocol: d.Protocol,
		Name:     d.Name,
		Address:  d.Parameter["address"],
		SlaveId:  byte(slaveid),
		State:    client.StateOffline,
		Commands: d.Commands,
	}, nil
}

func NewSerialClient(d *DeviceInfo) (*client.SerialClient, error) {
	if _, ok := d.Parameter["baudrate"]; !ok {
		return nil, fmt.Errorf("section 'parameter' should have field 'baudrate'")
	}
	baudrate, err := strconv.Atoi(d.Parameter["baudrate"])
	if err != nil {
		return nil, err
	}
	if _, ok := d.Parameter["path"]; !ok {
		return nil, fmt.Errorf("section 'parameter' should have field 'path'")
	}
	for _, c := range d.Commands {
		err = c.Check([]string{"bytes"})
		if err != nil {
			return nil, err
		}
	}
	return &client.SerialClient{
		Type:     d.Type,
		Protocol: d.Protocol,
		Name:     d.Name,
		Path:     d.Parameter["path"],
		BaudRate: baudrate,
		State:    client.StateOffline,
		Commands: d.Commands,
	}, nil
}

func NewTcpClient(d *DeviceInfo) (*client.TcpClient, error) {
	if _, ok := d.Parameter["address"]; !ok {
		return nil, fmt.Errorf("section 'parameter' should have field 'address'")
	}
	var err error
	for _, c := range d.Commands {
		err = c.Check([]string{"bytes"})
		if err != nil {
			return nil, err
		}
	}
	return &client.TcpClient{
		Type:     d.Type,
		Protocol: d.Protocol,
		Name:     d.Name,
		State:    client.StateOffline,
		Address:  d.Parameter["address"],
		Commands: d.Commands,
	}, nil
}

func NewSomeClient(d *DeviceInfo) (*client.VirtualClient, error) {
	return &client.VirtualClient{
		Type:     d.Type,
		Protocol: d.Protocol,
		Name:     d.Name,
		State:    client.StateOffline,
		Commands: d.Commands,
	}, nil
}
