package device

import (
	"encoding/json"
	"fmt"

	"adapter/db"
	"adapter/device/client"
	"adapter/device/command"
	"adapter/device/server"
)

var (
	TypeClient string = "client"
	TypeServer string = "server"
)

type DeviceInfo struct {
	Type      string            `json:"type"`     //"server" or "client"
	Protocol  string            `json:"protocol"` //"tcp", "modbus-tcp", "modbus-rtu", "iwsn", "serial", "virtual"
	Name      string            `json:"name"`
	Parameter map[string]string `json:"parameter"`
	Commands  []command.Command `json:"commands"`
}

func CreateDevice(d *DeviceInfo, insertdb bool) (err error) {

	switch d.Type {
	case TypeServer:
		ds, err := NewDeviceServer(d)
		if err != nil {
			return err
		}
		err = server.AddServerToList(ds)
		if err != nil {
			return err
		}
	case TypeClient:
		dc, err := NewDeviceClient(d)
		if err != nil {
			return err
		}
		err = client.AddClientToList(dc)
		if err != nil {
			return err
		}
	default:
		err = fmt.Errorf("Unknown support type: %s", d.Type)
	}
	if err == nil && insertdb {
		info, _ := json.Marshal(d)
		dbclient.DBStoreDevice(d.Type, d.Protocol, d.Name, info)
	}
	return
}

func GetAllDevices() []interface{} {
	dl := make([]interface{}, 0)
	for _, d := range client.GetClientList() {
		dl = append(dl, d)
	}
	for _, d := range server.GetServerList() {
		dl = append(dl, d)
	}
	return dl
}

func RemoveDevice(ttype string, name string) (err error) {
	switch ttype {
	case TypeClient:
		err = client.RemoveClientByName(name)
	case TypeServer:
		err = server.RemoveServerByName(name)
	default:
		err = fmt.Errorf("Unknown support type: %s", ttype)
	}
	if err == nil {
		dbclient.DBRemoveDevice(ttype, name)
	}
	return
}

func UpdateDevice(ttype string, name string, d *DeviceInfo) (err error) {
	switch ttype {
	case TypeClient:
		err = client.RemoveClientByName(name)
	case TypeServer:
		err = server.RemoveServerByName(name)
	default:
		err = fmt.Errorf("Unknown support type: %s", ttype)
	}
	if err == nil {
		dbclient.DBRemoveDevice(ttype, name)
		err = CreateDevice(d, true)
	}
	return
}

func ScheduleInit() {
	ds := dbclient.DBGetAllDevices()
	for _, s := range ds {
		var d DeviceInfo
		json.Unmarshal(s.Info, &d)
		CreateDevice(&d, false)
	}
}

func ScheduleStart() {
	ScheduleInit()
	go client.Schedule()
	go server.Schedule()
}

func ScheduleDestruct() {
	client.CleanupClientList()
	server.CleanupServerList()
}
