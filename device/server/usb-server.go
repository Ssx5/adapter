package server

import (
	"fmt"
	"strings"
	"time"

	"adapter/log"

	"adapter/nats"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

var (
	EtherTypeIpv4 string = "0x0800"
	EtherTypeIwsn string = "0x809a"
)

type UsbServer struct {
	handle       *pcap.Handle
	packetSource *gopacket.PacketSource
	mp           map[string]time.Time
	Type         string `json:"type"`
	Protocol     string `json:"protocol"`
	Name         string `json:"name"`
	Interface    string `json:"interface"`
	EthFilter    string `json:"ethfilter"`
	Timeout      string `json:"timeout"`
	State        string `json:"state"`
}

func (us *UsbServer) GetName() string {
	return us.Name
}

func (us *UsbServer) GetState() string {
	return us.State
}
func (us *UsbServer) Init() {
	var err error
	us.State = StateInit
	us.handle, err = pcap.OpenLive(us.Name, 65535, true, pcap.BlockForever)
	if err != nil {
		logclient.Log.Println(err)
		us.State = StateInitFailed
		return
	}
	logclient.Log.Printf("%s openlive success!", us.Interface)
	filter := fmt.Sprintf("ether[12:2] = %s", us.EthFilter)
	err = us.handle.SetBPFFilter(filter)
	if err != nil {
		logclient.Log.Fatalf("SetBPFFilter() error: %s\n", err)
	}
	logclient.Log.Printf("%s setfilter success!", us.Interface)
	us.packetSource = gopacket.NewPacketSource(us.handle, us.handle.LinkType())
	us.packetSource.NoCopy = true
	us.State = StateListening
	us.mp = make(map[string]time.Time)
	go us.checking()
	us.listening()
}

func (us *UsbServer) listening() {
	logclient.Log.Printf("%s start sniffing", us.Interface)
	for us.State == StateListening {
		packet := <-us.packetSource.Packets()
		id := strings.ToUpper(strings.Replace(packet.LinkLayer().LinkFlow().Src().String(), ":", "", -1))
		payload := packet.LinkLayer().LayerPayload()
		natsclient.Publish(us.getDataTopic(id), payload)
		t := time.Now()
		if _, ok := us.mp[id]; !ok {
			//send online message
			natsclient.Publish(natsclient.TopicEvent, us.getOnlineMessage(id))
		}
		us.mp[id] = t
	}
}

func (us *UsbServer) checking() {
	duration, _ := time.ParseDuration(us.Timeout)
	ticker := time.NewTicker(duration / 10)
	for us.State == StateListening {
		select {
		case <-ticker.C:
			t := time.Now()
			for id, v := range us.mp {
				if t.Unix()-v.Unix() >= int64(duration) {
					//send offline message
					natsclient.Publish(natsclient.TopicEvent, us.getOfflineMessage(id))
					delete(us.mp, id)
				}
			}
		}
	}
}

func (us *UsbServer) Close() {
	us.handle.Close()
	us.State = StateDeleted
}

func (us *UsbServer) getOnlineMessage(id string) []byte {
	return []byte(fmt.Sprintf(`{"type":"%s", "id":"%s", "event":"online"}`, us.Protocol, id))
}
func (us *UsbServer) getOfflineMessage(id string) []byte {
	return []byte(fmt.Sprintf(`{"type":"%s", "id":"%s", "event":"offline"}`, us.Protocol, id))
}

func (us *UsbServer) getDataTopic(id string) string {
	return fmt.Sprintf("dev.data.%s.%s", us.Name, id)
}
