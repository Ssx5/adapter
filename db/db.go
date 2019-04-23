package dbclient

import (
	"adapter/log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var db *gorm.DB

func NewInstance(path string) {
	var err error
	if db != nil {
		db.Close()
	}
	db, err = gorm.Open("sqlite3", path)
	if err != nil {
		logclient.Log.Println(err)
	}
	db.AutoMigrate(&DeviceInfoSeries{})
}

type DeviceInfoSeries struct {
	gorm.Model
	Type     string
	Protocol string
	Name     string
	Info     []byte
}

func DBStoreDevice(ttype string, protocol string, name string, info []byte) {
	ds := &DeviceInfoSeries{
		Type:     ttype,
		Protocol: protocol,
		Name:     name,
		Info:     info,
	}
	db.Create(ds)
}

func DBGetAllDevices() []DeviceInfoSeries {
	var ds []DeviceInfoSeries
	db.Find(&ds)
	return ds
}

func DBRemoveDevice(ttype string, name string) {
	db.Unscoped().Where("type = ? AND name = ?", ttype, name).Delete(&DeviceInfoSeries{})
}
