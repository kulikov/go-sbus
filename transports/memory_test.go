package transports

import (
	"testing"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/kulikov/go-sbus"
)

type User struct {
	Name string
	Age  int
}

func TestInMemory_Pub(t *testing.T) {
	memory := NewInMemory(logrus.New().WithField("logger", "sbus"))

	mark := 0

	memory.Sub("get-users", func(msg sbus.Message, log *logrus.Entry) error {
		log.Infof("Received one %v", msg)
		mark += 1
		return nil
	})

	memory.Sub("get-users", func(msg sbus.Message, log *logrus.Entry) error {
		log.Infof("Received two %v", msg)
		mark += 1
		return nil
	})

	memory.Pub(&sbus.Message{
		Subject: "get-users",
		Data:    map[string]string{"id": "123"},
		Meta:    sbus.Meta{"timestamp": 131312424},
	})

	time.Sleep(5 * time.Millisecond)

	if mark != 2 {
		t.Fail()
	}
}

func TestInMemory_SubOnce(t *testing.T) {
	memory := NewInMemory(logrus.New().WithField("logger", "sbus"))

	mark := 0

	memory.SubOnce("do-once", func(msg sbus.Message, log *logrus.Entry) error {
		log.Infof("Received once %v", msg)
		mark += 1
		return nil
	})

	memory.Pub(&sbus.Message{Subject: "do-once"})
	memory.Pub(&sbus.Message{Subject: "do-once"})
	memory.Pub(&sbus.Message{Subject: "do-once"})

	time.Sleep(5 * time.Millisecond)

	if mark != 1 {
		t.Fail()
	}
}

func TestInMemory_SubOnceMixed(t *testing.T) {
	memory := NewInMemory(logrus.New().WithField("logger", "sbus"))

	mark := 0

	memory.SubOnce("do-any", func(msg sbus.Message, log *logrus.Entry) error {
		log.Infof("Received once %v", msg)
		mark += 1
		return nil
	})

	memory.Sub("do-any", func(msg sbus.Message, log *logrus.Entry) error {
		log.Infof("Received %v", msg)
		mark += 1
		return nil
	})

	memory.Pub(&sbus.Message{Subject: "do-any"})
	memory.Pub(&sbus.Message{Subject: "do-any"})
	memory.Pub(&sbus.Message{Subject: "do-any"})

	time.Sleep(5 * time.Millisecond)

	if mark != 4 {
		t.Fail()
	}
}
