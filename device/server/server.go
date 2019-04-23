package server

import (
	"fmt"
	"sync"
	"time"

	"adapter/log"
)

type DeviceServer interface {
	Init()
	GetName() string
	GetState() string
	Close()
}

type DeviceServerList struct {
	list map[string]DeviceServer
	lock *sync.Mutex
}

var deviceServerList DeviceServerList

func init() {
	deviceServerList.list = make(map[string]DeviceServer)
	deviceServerList.lock = &sync.Mutex{}

}

func (dl *DeviceServerList) Getlist() map[string]DeviceServer {
	return dl.list
}

func (dl *DeviceServerList) AddServer(ds DeviceServer) error {
	if _, ok := dl.list[ds.GetName()]; ok {
		return fmt.Errorf("device %s already exists!", ds.GetName())
	}
	dl.list[ds.GetName()] = ds
	return nil
}

func (dl *DeviceServerList) FindServerByName(name string) (DeviceServer, error) {
	if _, ok := dl.list[name]; !ok {
		return nil, fmt.Errorf("device %s doesn't exists!", name)
	}
	return dl.list[name], nil
}

func (dl *DeviceServerList) RemoveServerByName(name string) error {
	if _, ok := dl.list[name]; !ok {
		return fmt.Errorf("device %s doesn't exists!", name)
	}
	go dl.list[name].Close()
	delete(dl.list, name)
	return nil
}

func (dl *DeviceServerList) Cleanup() {
	for _, s := range dl.list {
		if s.GetState() == StateListening {
			s.Close()
		}
	}
	dl.list = nil
}

func Schedule() {
	logclient.Log.Printf("server schedule starts!")
	ticker := time.NewTicker(time.Second)
	for {
		select {
		case <-ticker.C:
			ServerProcess(deviceServerList)
		}
	}
}

func ServerProcess(dl DeviceServerList) {
	for _, d := range dl.list {
		switch d.GetState() {
		case StateOffline, StateInitFailed:
			go d.Init()
		case StateInit, StateListening, StateDeleted:
			//do nonthing
		}
	}
}

func GetServerList() map[string]DeviceServer {
	return deviceServerList.Getlist()
}

func AddServerToList(ds DeviceServer) error {
	return deviceServerList.AddServer(ds)
}

func FindServerByName(name string) (DeviceServer, error) {
	return deviceServerList.FindServerByName(name)
}

func RemoveServerByName(name string) error {
	return deviceServerList.RemoveServerByName(name)
}

func CleanupServerList() {
	deviceServerList.Cleanup()
}
