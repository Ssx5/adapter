package client

import (
	"fmt"
	"strconv"
	"time"

	"github.com/goburrow/modbus"
	"adapter/device/command"
	"adapter/log"
	"adapter/nats"
)

type ModbusRtuClient struct {
	handler  *modbus.RTUClientHandler
	client   modbus.Client
	Type     string            `json:"type"`
	Protocol string            `json:"protocol"`
	Name     string            `json:"name"`
	Path     string            `json:"path"`
	BaudRate int               `json:"baudrate"`
	DataBits int               `json:"databits"`
	Parity   string            `json:"parity"`
	StopBits int               `json:"stopbits"`
	SlaveId  byte              `json:"slaveid"`
	State    string            `json:"state"`
	Commands []command.Command `json:"commands"`
}

func (mrc *ModbusRtuClient) GetName() string {
	return mrc.Name
}

func (mrc *ModbusRtuClient) GetState() string {
	return mrc.State
}

func (mrc *ModbusRtuClient) Connect() {
	mrc.handler = modbus.NewRTUClientHandler(mrc.Path)
	mrc.handler.BaudRate = mrc.BaudRate
	mrc.handler.DataBits = mrc.DataBits
	mrc.handler.Parity = mrc.Parity
	mrc.handler.StopBits = mrc.StopBits
	mrc.handler.SlaveId = mrc.SlaveId
	mrc.handler.Timeout = 5 * time.Second

	mrc.setState(StateConnecting)
	err := mrc.handler.Connect()
	if err != nil {
		logclient.Log.Println(err)
		mrc.setState(StateConnectFailed)
		return
	}
	mrc.client = modbus.NewClient(mrc.handler)
	mrc.setState(StateOnline)
}

func (mrc *ModbusRtuClient) GetCommands() []command.Command {
	return mrc.Commands
}

func (mrc *ModbusRtuClient) GetCommandByName(Name string) (c *command.Command, err error) {
	for k, v := range mrc.Commands {
		if v.Name == Name {
			c = &mrc.Commands[k]
			break
		}
	}
	if c == nil {
		err = fmt.Errorf("%s doesn't have command %s\n", mrc.Name, Name)
	}
	return c, err
}

func (mrc *ModbusRtuClient) ExecCommand(c command.Command) {
	funCode, _ := strconv.Atoi(c.Attribution["funCode"].(string))
	startAddr, _ := strconv.Atoi(c.Attribution["startAddr"].(string))
	quantity, _ := strconv.Atoi(c.Attribution["quantity"].(string))
	CommandFunc := mrc.getReadFunc(funCode)
	results, err := CommandFunc(uint16(startAddr), uint16(quantity))
	if err != nil {
		logclient.Log.Printf("modbus read error %s\n", err)
		return
	}
	natsclient.Publish(mrc.getDataTopic(c.Name), results)
}

func (mrc *ModbusRtuClient) Close() {
	if mrc.State == StateOnline {
		mrc.handler.Close()
	}
	mrc.setState(StateOffline)
}

/*
Private
*/
func (mrc *ModbusRtuClient) setState(s string) {
	if mrc.State == s {
		return
	}
	if mrc.State != StateOnline && s == StateOnline {
		natsclient.Publish(natsclient.TopicEvent, mrc.getOnlineMessage())
	} else if mrc.State == StateOnline {
		natsclient.Publish(natsclient.TopicEvent, mrc.getOfflineMessage())
	}
	mrc.State = s
}

func (mrc *ModbusRtuClient) getOnlineMessage() []byte {
	return []byte(fmt.Sprintf(`{"type":"%s", "id":"%s", "event":"online"}`, mrc.Protocol, mrc.Name))
}

func (mrc *ModbusRtuClient) getOfflineMessage() []byte {
	return []byte(fmt.Sprintf(`{"type":"%s", "id":"%s", "event":"offline"}`, mrc.Protocol, mrc.Name))
}

func (mrc *ModbusRtuClient) getDataTopic(cname string) string {
	return fmt.Sprintf("%s.%s.%s", natsclient.TopicData, cname, mrc.Path)
}

type ModbusRtuReadFunc func(uint16, uint16) ([]byte, error)

func (mrc *ModbusRtuClient) getReadFunc(funCode int) ModbusRtuReadFunc {
	var CommandFunc ModbusRtuReadFunc
	switch funCode {
	case 0x01:
		CommandFunc = mrc.client.ReadCoils
	case 0x02:
		CommandFunc = mrc.client.ReadDiscreteInputs
	case 0x03:
		CommandFunc = mrc.client.ReadHoldingRegisters
	case 0x04:
		CommandFunc = mrc.client.ReadInputRegisters
	}
	return CommandFunc
}
