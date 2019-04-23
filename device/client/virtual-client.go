package client

import (
	"fmt"

	"adapter/nats"

	"adapter/device/command"
)

type VirtualClient struct {
	count    int
	Type     string            `json:"type"`
	Protocol string            `json:"protocol"`
	Name     string            `json:"name"`
	State    string            `json:"state"`
	Commands []command.Command `json:"commands"`
}

func (s *VirtualClient) GetName() string {
	return s.Name
}

func (s *VirtualClient) Connect() {
	s.count = 0
	s.setState(StateOnline)
}

func (s *VirtualClient) GetState() string {
	return s.State
}

func (s *VirtualClient) GetCommands() []command.Command {
	return s.Commands
}

func (s *VirtualClient) ExecCommand(c command.Command) {
	s.count++
	results := fmt.Sprintf("%d %s VirtualClient exec commands %s Attribution: %v Period: %v\n", s.count, s.Name, c.Name, c.Attribution, c.Period)
	natsclient.Publish(s.getDataTopic(c.Name), []byte(results))
}

func (s *VirtualClient) GetCommandByName(Name string) (c *command.Command, err error) {
	for k, v := range s.Commands {
		if v.Name == Name {
			c = &s.Commands[k]
			break
		}
	}
	if c == nil {
		err = fmt.Errorf("%s doesn't have command %s\n", s.Name, Name)
	}
	return c, err
}

func (s *VirtualClient) Close() {
	s.setState(StateOffline)
}

/*
Private
*/
func (sc *VirtualClient) setState(s string) {
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

func (s *VirtualClient) getOnlineMessage() []byte {
	return []byte(fmt.Sprintf(`{"type":"%s", "id":"%s", "event":"online"}`, s.Protocol, s.Name))
}

func (s *VirtualClient) getOfflineMessage() []byte {
	return []byte(fmt.Sprintf(`{"type":"%s", "id":"%s", "event":"offline"}`, s.Protocol, s.Name))
}
func (s *VirtualClient) getDataTopic(cname string) string {
	return fmt.Sprintf("%s.%s.%s", natsclient.TopicData, cname, s.Name)
}
