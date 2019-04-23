package client

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/goburrow/modbus"
	"adapter/device/command"
	"adapter/log"
	"adapter/nats"
)

type ModbusTcpClient struct {
	handler  *modbus.TCPClientHandler
	client   modbus.Client
	Type     string            `json:"type"`
	Protocol string            `json:"protocol"`
	Name     string            `json:"name"`
	Address  string            `josn:"address"`
	SlaveId  byte              `json:"slaveid"`
	State    string            `json:"state"`
	Commands []command.Command `json:"commands"`
}

func (mtc *ModbusTcpClient) GetName() string {
	return mtc.Name
}

func (mtc *ModbusTcpClient) GetState() string {
	return mtc.State
}

func (mtc *ModbusTcpClient) Connect() {

	mtc.handler = modbus.NewTCPClientHandler(mtc.Address)
	mtc.handler.Timeout = 10 * time.Second
	mtc.handler.SlaveId = mtc.SlaveId
	mtc.setState(StateConnecting)
	err := mtc.handler.Connect()
	if err != nil {
		logclient.Log.Println(err)
		mtc.setState(StateConnectFailed)
		return
	}
	mtc.client = modbus.NewClient(mtc.handler)
	mtc.setState(StateOnline)
}

func (mtc *ModbusTcpClient) GetCommands() []command.Command {
	return mtc.Commands
}

func (mtc *ModbusTcpClient) GetCommandByName(Name string) (c *command.Command, err error) {
	for k, v := range mtc.Commands {
		if v.Name == Name {
			c = &mtc.Commands[k]
			break
		}
	}
	if c == nil {
		err = fmt.Errorf("%s doesn't have command %s\n", mtc.Name, Name)
	}
	return c, err
}

func (mtc *ModbusTcpClient) ExecCommand(c command.Command) {
	funCode, _ := strconv.Atoi(c.Attribution["funCode"].(string))
	startAddr, _ := strconv.Atoi(c.Attribution["startAddr"].(string))
	quantity, _ := strconv.Atoi(c.Attribution["quantity"].(string))
	CommandFunc := mtc.getReadFunc(funCode)
	results, err := CommandFunc(uint16(startAddr), uint16(quantity))
	if err != nil {
		logclient.Log.Printf("modbus read error %s\n", err)
		return
	}
	natsclient.Publish(mtc.getDateTopic(c.Name), results)
}

func (mtc *ModbusTcpClient) Close() {
	if mtc.State == StateOnline {
		mtc.handler.Close()
	}
	mtc.setState(StateOffline)
	/*
		nats send offline packet here
	*/
}

func (mtc *ModbusTcpClient) setState(s string) {
	if mtc.State == s {
		return
	}
	if mtc.State != StateOnline && s == StateOnline {
		natsclient.Publish(natsclient.TopicEvent, mtc.getOnlineMessage())
	} else if mtc.State == StateOnline {
		natsclient.Publish(natsclient.TopicEvent, mtc.getOfflineMessage())
	}
	mtc.State = s
}

func (mtc *ModbusTcpClient) getOnlineMessage() []byte {
	return []byte(fmt.Sprintf(`{"type":"%s", "id":"%s", "event":"online"}`, mtc.Protocol, mtc.Name))
}

func (mtc *ModbusTcpClient) getOfflineMessage() []byte {
	return []byte(fmt.Sprintf(`{"type":"%s", "id":"%s", "event":"offline"}`, mtc.Protocol, mtc.Name))
}

func (mtc *ModbusTcpClient) getDateTopic(cname string) string {
	return fmt.Sprintf("%s.%s.%s", natsclient.TopicData, cname, strings.Replace(mtc.Address, ".", "_", -1))
}

type ModbusTcpReadFunc func(uint16, uint16) ([]byte, error)

func (mtc *ModbusTcpClient) getReadFunc(funCode int) ModbusTcpReadFunc {
	var CommandFunc ModbusTcpReadFunc
	switch funCode {
	case 0x01:
		CommandFunc = mtc.client.ReadCoils
	case 0x02:
		CommandFunc = mtc.client.ReadDiscreteInputs
	case 0x03:
		CommandFunc = mtc.client.ReadHoldingRegisters
	case 0x04:
		CommandFunc = mtc.client.ReadInputRegisters
	}
	return CommandFunc
}
