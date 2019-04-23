package config

import (
	"flag"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"adapter/log"
	"adapter/nats"
)

type AdapterConfig struct {
	gorm.Model `json:"-"`
	NatsUrl    string `json:"natsurl"`
	DeviceDB   string `json:"devicedb"`
	RestPort   int    `json:"restport"`
}

var configDBPath string
var configdb *gorm.DB
var globalConfig AdapterConfig

func GetGloablConfig() AdapterConfig {
	return globalConfig
}

func AdapterConfigInit() {
	flag.StringVar(&configDBPath, "db", "", "Path to the database file of configuration")
	flag.Parse()
	if configDBPath == "" {
		flag.Usage()
		os.Exit(1)
	}
	var err error
	configdb, err = gorm.Open("sqlite3", configDBPath)
	if err != nil {
		logclient.Log.Fatalf(err.Error())
	}
	configdb.AutoMigrate(&AdapterConfig{})
	configdb.First(&globalConfig)
	if globalConfig.ID != 1 {
		globalConfig = AdapterConfig{
			NatsUrl:  "localhost:4222",
			DeviceDB: configDBPath,
			RestPort: 9997,
		}
		configdb.Create(&globalConfig)
	}
	logclient.Log.Printf("globalConfig: natsurl: %s, devicedb: %s, restport: %d\n", globalConfig.NatsUrl, globalConfig.DeviceDB, globalConfig.RestPort)
}

func SetNewConfig(c AdapterConfig) {
	if c.NatsUrl != globalConfig.NatsUrl {
		natsclient.NewInstance(c.NatsUrl)
		configdb.First(&globalConfig)
		globalConfig.NatsUrl = c.NatsUrl
		configdb.Save(&globalConfig)
	}

	if c.RestPort != globalConfig.RestPort {
		configdb.First(&globalConfig)
		globalConfig.RestPort = c.RestPort
		configdb.Save(&globalConfig)
	}
}
