package natsclient

import (
	"github.com/nats-io/go-nats"
	"adapter/log"
)

var (
	TopicEvent string = "dev.event"
	TopicData  string = "dev.data"
)

var ns *nats.Conn

func NewInstance(natsurl string) {
	var err error
	if ns != nil {
		ns.Flush()
		ns.Close()
	}
	ns, err = nats.Connect(natsurl)
	if err != nil {
		logclient.Log.Fatalf("nats cannot connect %s: %s", natsurl, err)
	}
	logclient.Log.Printf("nats connect %s success!", natsurl)
}

func Publish(topic string, msg []byte) {
	ns.Publish(topic, msg)
	ns.Flush()
}
