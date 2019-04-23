package server

import (
	"fmt"
	"net"
	"strings"
	"time"

	"adapter/log"

	"adapter/nats"
)

type TcpServer struct {
	listener net.Listener
	mp       map[string]time.Time
	Type     string `json:"type"`
	Protocol string `json:"protocol"`
	Name     string `json:"name"`
	Address  string `json:"address"`
	State    string `json:"state"`
	Timeout  string `json:"timeout"`
}

func (ts *TcpServer) GetName() string {
	return ts.Name
}

func (ts *TcpServer) GetState() string {
	return ts.State
}

func (ts *TcpServer) Init() {
	var err error
	ts.State = StateInit
	ts.listener, err = net.Listen("tcp", ts.Address)

	if err != nil {
		logclient.Log.Println(err)
		ts.State = StateInitFailed
	}
	ts.State = StateListening
	logclient.Log.Printf("%s is listening at: %s", ts.Name, ts.Address)
	ts.mp = make(map[string]time.Time)
	go ts.checking()
	ts.listening()
}

func (ts *TcpServer) listening() {
	for ts.State == StateListening {
		conn, err := ts.listener.Accept()
		if err != nil {
			logclient.Log.Printf("%s", err)
			break
		}
		logclient.Log.Printf("%s accepts a client: %s", ts.Name, conn.RemoteAddr())
		// start a new goroutine to handle the new connection
		go func(ts *TcpServer, conn net.Conn) {
			defer conn.Close()
			for ts.State == StateListening {
				var results = make([]byte, 1024)
				n, err := conn.Read(results)
				if err != nil {
					logclient.Log.Println("conn read error:", err)
					return
				}
				logclient.Log.Printf("%s receive: %s", ts.Name, string(results[:n]))
				id := strings.Replace(fmt.Sprintf("%s", conn.RemoteAddr()), ".", "_", -1)
				t := time.Now()
				if _, ok := ts.mp[id]; !ok {
					//send online message
					natsclient.Publish(natsclient.TopicEvent, ts.getOnlineMessage(id))
					logclient.Log.Println(id, "online")
				}
				ts.mp[id] = t
				// upload data
				natsclient.Publish(ts.getDataTopic(id), results[:n])
			}
		}(ts, conn)
	}
}

func (ts *TcpServer) checking() {
	duration, _ := time.ParseDuration(ts.Timeout)
	ticker := time.NewTicker(duration / 10)
	for ts.State == StateListening {
		select {
		case <-ticker.C:
			t := time.Now()
			for id, v := range ts.mp {
				if t.Unix()-v.Unix() >= int64(duration)/int64(time.Second) {
					//send offline message
					natsclient.Publish(natsclient.TopicEvent, ts.getOfflineMessage(id))
					delete(ts.mp, id)
				}
			}
		}
	}
}

func (ts *TcpServer) Close() {
	ts.listener.Close()
	ts.State = StateDeleted
}

func (ts *TcpServer) getOnlineMessage(id string) []byte {
	return []byte(fmt.Sprintf(`{"type":"%s", "id":"%s", "event":"online"}`, ts.Protocol, id))
}

func (ts *TcpServer) getOfflineMessage(id string) []byte {
	return []byte(fmt.Sprintf(`{"type":"%s", "id":"%s", "event":"offline"}`, ts.Protocol, id))
}

func (ts *TcpServer) getDataTopic(id string) string {
	return fmt.Sprintf("dev.data.%s.%s", ts.Name, id)
}
