package natsclient

import (
	"strconv"
	"testing"
)

func TestPublish(t *testing.T) {
	NewInstance("127.0.0.1:4222")
	for i := 0; i < 10; i++ {
		Publish("test."+strconv.Itoa(i)+"y", []byte("hello world"))
	}
}
