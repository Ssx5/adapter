package client

import (
	"fmt"
	"sync"
	"time"

	"adapter/device/command"
	"adapter/log"
)

type DeviceClient interface {
	GetName() string
	Connect()
	GetState() string
	GetCommands() []command.Command
	GetCommandByName(Name string) (*command.Command, error)
	ExecCommand(command.Command)
	Close()
}

func Schedule() {
	ticker := time.NewTicker(time.Millisecond * 100)
	logclient.Log.Printf("client schedule starts!")
	for {
		select {
		case <-ticker.C:
			ClientProcess(deviceClientList)
		}
	}
}

func ClientProcess(dl DeviceList) {
	for _, dc := range dl.GetList() {
		switch dc.GetState() {
		case StateOnline:
			commands := dc.GetCommands()
			for i := range commands {
				t := time.Now()
				period, _ := time.ParseDuration(commands[i].Period)
				if period > 0 && t.UnixNano()/10000000 >= commands[i].NextTime.UnixNano()/10000000 {
					go func(dc DeviceClient, i int) {
						dc.ExecCommand(commands[i])
						commands[i].NextTime = t.Add(period)
					}(dc, i)
				}
			}
		case StateOffline, StateConnectFailed:
			go dc.Connect()
		case StateConnecting:
			//do nothing
		}
	}
}

/////////////////////////////////////////////////////////////////////////////////////////////

type DeviceList struct {
	list map[string]DeviceClient
	lock *sync.Mutex
}

var deviceClientList DeviceList

func init() {
	deviceClientList.list = make(map[string]DeviceClient)
	deviceClientList.lock = &sync.Mutex{}
}

func (dl *DeviceList) GetList() map[string]DeviceClient {
	return dl.list
}

func (dl *DeviceList) AddClient(c DeviceClient) error {
	if _, ok := dl.list[c.GetName()]; ok {
		return fmt.Errorf("Client %s already exists", c.GetName())
	}
	dl.list[c.GetName()] = c
	return nil
}

func (dl *DeviceList) FindClientByName(name string) (DeviceClient, error) {
	c, ok := dl.list[name]
	if !ok {
		return nil, fmt.Errorf("Client %s doesn't exist", name)
	}
	return c, nil
}

func (dl *DeviceList) RemoveClientByName(name string) error {
	if _, ok := dl.list[name]; !ok {
		return fmt.Errorf("Client %s doesn't exist", name)
	}
	go dl.list[name].Close()
	delete(dl.list, name)
	return nil
}

func (dl *DeviceList) Cleanup() {
	for _, c := range dl.list {
		if c.GetState() == StateOnline {
			c.Close()
		}
	}
	dl.list = nil
}

func AddClientToList(c DeviceClient) error {
	return deviceClientList.AddClient(c)
}

func GetClientList() map[string]DeviceClient {
	return deviceClientList.GetList()
}

func FindClientByName(name string) (DeviceClient, error) {
	return deviceClientList.FindClientByName(name)
}

func RemoveClientByName(name string) error {
	return deviceClientList.RemoveClientByName(name)
}

func CleanupClientList() {
	deviceClientList.Cleanup()
}
