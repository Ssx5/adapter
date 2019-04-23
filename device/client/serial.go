package client

import (
	"fmt"
	"sync"

	"adapter/log"

	"adapter/device/command"
	"adapter/nats"
	"github.com/tarm/serial"
)

type SerialClient struct {
	s        *serial.Port
	lock     sync.Mutex
	Type     string            `json:"type"`
	Protocol string            `json:"protocol"`
	Name     string            `json:"name"`
	Path     string            `json:"path"`
	BaudRate int               `json:"baudrate"`
	State    string            `json:"state"`
	Commands []command.Command `json:"commands"`
}

func (sc *SerialClient) GetName() string {
	return sc.Name
}

func (sc *SerialClient) Connect() {
	c := &serial.Config{Name: sc.Path, Baud: sc.BaudRate}
	var err error
	sc.setState(StateConnecting)
	sc.s, err = serial.OpenPort(c)
	if err != nil {
		logclient.Log.Println(err)
		sc.setState(StateConnectFailed)
		return
	}
	sc.setState(StateOnline)
}

func (sc *SerialClient) GetState() string {
	return sc.State
}

func (sc *SerialClient) GetCommands() []command.Command {
	return sc.Commands
}

func (sc *SerialClient) GetCommandByName(Name string) (c *command.Command, err error) {
	for k, v := range sc.Commands {
		if v.Name == Name {
			c = &sc.Commands[k]
			break
		}
	}
	if c == nil {
		err = fmt.Errorf("%s doesn't have command %s\n", sc.Name, Name)
	}
	return c, err
}
func (sc *SerialClient) ExecCommand(c command.Command) {
	stream := c.Attribution["bytes"]
	sc.lock.Lock()
	_, err := sc.s.Write(stream.([]byte))
	if err != nil {
		logclient.Log.Println(err)
		return
	}
	results := make([]byte, 128)
	n, err := sc.s.Read(results)
	if err != nil {
		logclient.Log.Println(err)
		return
	}
	sc.lock.Unlock()
	natsclient.Publish(sc.getDataTopic(c.Name), results[:n])
}

func (sc *SerialClient) Close() {
	if sc.State == StateOnline {
		sc.s.Close()
	}
	sc.setState(StateOffline)
}

func (sc *SerialClient) setState(s string) {
	if sc.State == s {
		return
	}
	if sc.State != StateOnline && s == StateOnline {
		natsclient.Publish(natsclient.TopicEvent, sc.getOnlineMessage())
	} else if sc.State == StateOnline {
		natsclient.Publish(natsclient.TopicEvent, sc.getOfflineMessage())
	}
	sc.State = s
}

func (sc *SerialClient) getOnlineMessage() []byte {
	return []byte(fmt.Sprintf(`{"type":"%s", "id":"%s", "event":"online"}`, sc.Protocol, sc.Name))
}

func (sc *SerialClient) getOfflineMessage() []byte {
	return []byte(fmt.Sprintf(`{"type":"%s", "id":"%s", "event":"offline"}`, sc.Protocol, sc.Name))
}

func (sc *SerialClient) getDataTopic(cname string) string {
	return fmt.Sprintf("%s.%s.%s", natsclient.TopicData, cname, sc.Path)
}
