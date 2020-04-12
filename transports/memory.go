package transports

import (
	"sync"

	"github.com/Sirupsen/logrus"

	"github.com/kulikov/go-sbus"
)

func NewInMemory(log *logrus.Entry) *InMemory {
	return &InMemory{
		log:       log,
		listeners: make(map[string][]sbus.MessageHandler, 0),
		once:      make(map[string][]sbus.MessageHandler, 0),
	}
}

type InMemory struct {
	log       *logrus.Entry
	listeners map[string][]sbus.MessageHandler
	once      map[string][]sbus.MessageHandler
	sync.RWMutex
}

func (t *InMemory) Sub(subject string, handler sbus.MessageHandler) error {
	t.Lock()
	defer t.Unlock()

	exists, ok := t.listeners[subject]
	if !ok {
		exists = make([]sbus.MessageHandler, 0)
	}

	t.listeners[subject] = append(exists, handler)
	return nil
}

func (t *InMemory) SubOnce(subject string, handler sbus.MessageHandler) error {
	t.Lock()
	defer t.Unlock()

	exists, ok := t.once[subject]
	if !ok {
		exists = make([]sbus.MessageHandler, 0)
	}

	t.once[subject] = append(exists, handler)
	return nil
}

func (t *InMemory) Pub(msg *sbus.Message) error {
	t.RLock()
	listeners := t.listeners[msg.Subject]
	_, onceOk := t.once[msg.Subject]
	t.RUnlock()

	if onceOk {
		t.Lock()
		if once, ok := t.once[msg.Subject]; ok {
			delete(t.once, msg.Subject)
			t.Unlock()
			listeners = append(listeners, once...)
		} else {
			t.Unlock()
		}
	}

	t.log.Debugf("Publish %v ~~~> %s", msg.Subject, msg.Data)

	for _, handler := range listeners {
		go func(msg *sbus.Message, hdr sbus.MessageHandler) {
			if err := hdr(*msg); err != nil {
				t.log.WithError(err).Errorf("Error on handle message %v ~~~> %s", msg.Subject, msg.Data)
			}
		}(msg, handler)
	}

	return nil
}
