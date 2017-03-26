package transports

import (
	"fmt"
	"testing"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/kulikov/go-sbus"
)

type User struct {
	Name string
	Age  int
}

func TestInMemory_Pub(t *testing.T) {
	memory := NewInMemory(log.New().WithField("logger", "sbus"))

	mark := 0

	memory.Sub("get-users", func(msg sbus.Message) error {
		log.Infof("Received one %v", msg)

		u := &User{}
		msg.Unmarshal(u)
		log.Infof("User: %v", u)

		if fmt.Sprintf("%v", u) != fmt.Sprintf("%v", &User{"Dima", 31}) {
			t.Fail()
		}

		mark += 1
		return nil
	})

	memory.Sub("get-users", func(msg sbus.Message) error {
		log.Infof("Received two %v", msg)
		mark += 1
		return nil
	})

	memory.Pub(&sbus.Message{
		Subject: "get-users",
		Data:    sbus.Marshal(User{"Dima", 31}),
		Meta:    sbus.Meta{"timestamp": 131312424},
	})

	time.Sleep(5 * time.Millisecond)

	if mark != 2 {
		t.Fail()
	}
}

func TestInMemory_SubOnce(t *testing.T) {
	memory := NewInMemory(log.New().WithField("logger", "sbus"))

	mark := 0

	memory.SubOnce("do-once", func(msg sbus.Message) error {
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
	memory := NewInMemory(log.New().WithField("logger", "sbus"))

	mark := 0

	memory.SubOnce("do-any", func(msg sbus.Message) error {
		log.Infof("Received once %v", msg)
		mark += 1
		return nil
	})

	memory.Sub("do-any", func(msg sbus.Message) error {
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
