package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"adapter/log"

	"adapter/config"
	"adapter/db"
	"adapter/device"
	"adapter/nats"
	"adapter/rest"
)

func main() {

	config.AdapterConfigInit()
	natsclient.NewInstance(config.GetGloablConfig().NatsUrl)
	dbclient.NewInstance(config.GetGloablConfig().DeviceDB)

	device.ScheduleStart()

	errs := make(chan error, 2)
	rest.StartHttpServer(errs, config.GetGloablConfig().RestPort)
	listenForInterrupt(errs)
	<-errs
	device.ScheduleDestruct()
	logclient.Log.Printf("Adapter has been interrupt!")
}

func listenForInterrupt(errChan chan error) {
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errChan <- fmt.Errorf("%s", <-c)
	}()
}
