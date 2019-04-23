package client

import (
	"fmt"
	"strings"
	"sync"

	"net"

	"adapter/device/command"
	"adapter/log"
	"adapter/nats"
)

type TcpClient struct {
	conn     net.Conn
	lock     sync.Mutex
	Type     string            `json:"type"`
	Protocol string            `json:"protocol"`
	Name     string            `json:"name"`
	State    string            `json:"state"`
	Address  string            `json:"address"`
	Commands []command.Command `json:"commands"`
}

func (tc *TcpClient) GetName() string {
	return tc.Name
}

func (tc *TcpClient) Connect() {

	tc.setState(StateConnecting)
	var err error
	tc.conn, err = net.Dial("tcp", tc.Address)
	if err != nil {
		logclient.Log.Printf("%s", err)
		tc.setState(StateConnectFailed)
		return
	}
	tc.setState(StateOnline)
	return
}

func (tc *TcpClient) GetState() string {
	return tc.State
}

func (tc *TcpClient) GetCommands() []command.Command {
	return tc.Commands
}

func (tc *TcpClient) GetCommandByName(Name string) (c *command.Command, err error) {
	for k, v := range tc.Commands {
		if v.Name == Name {
			c = &tc.Commands[k]
			break
		}
	}
	if c == nil {
		err = fmt.Errorf("%s doesn't have command %s\n", tc.Name, Name)
	}
	return c, err
}
func (tc *TcpClient) ExecCommand(c command.Command) {
	stream := c.Attribution["bytes"]
	tc.lock.Lock()
	_, err := tc.conn.Write(stream.([]byte))
	if err != nil {
		logclient.Log.Println(err)
		return
	}
	results := make([]byte, 128)
	n, err := tc.conn.Read(results)
	if err != nil {
		logclient.Log.Println(err)
		return
	}
	tc.lock.Unlock()
	//nats upload here
	natsclient.Publish(tc.getDateTopic(c.Name), results[:n])
}

func (tc *TcpClient) Close() {
	if tc.State == StateOnline {
		tc.conn.Close()
	}
	tc.setState(StateOffline)
}

/*
Private
*/

func (tc *TcpClient) setState(s string) {
	if tc.State == s {
		return
	}
	if tc.State != StateOnline && s == StateOnline {
		natsclient.Publish(natsclient.TopicEvent, tc.getOnlineMessage())
	} else if tc.State == StateOnline {
		natsclient.Publish(natsclient.TopicEvent, tc.getOfflineMessage())
	}
	tc.State = s
}

func (tc *TcpClient) getOnlineMessage() []byte {
	return []byte(fmt.Sprintf(`{"type":"%s", "id":"%s", "event":"online"}`, tc.Protocol, tc.Name))
}

func (tc *TcpClient) getOfflineMessage() []byte {
	return []byte(fmt.Sprintf(`{"type":"%s", "id":"%s", "event":"offline"}`, tc.Protocol, tc.Name))
}

func (tc *TcpClient) getDateTopic(cname string) string {
	return fmt.Sprintf("%s.%s.%s", natsclient.TopicData, cname, strings.Replace(tc.Address, ".", "_", -1))
}
