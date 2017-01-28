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

type handlerItem struct {
	handler sbus.MessageHandler
	once    bool
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
	exists := t.listeners[msg.Subject]
	_, onceOk := t.once[msg.Subject]
	t.RUnlock()

	if onceOk {
		t.Lock()
		if once, ok := t.once[msg.Subject]; ok {
			delete(t.once, msg.Subject)
			t.Unlock()
			exists = append(exists, once...)
		} else {
			t.Unlock()
		}
	}

	if len(exists) > 0 {
		for _, handler := range exists {
			go func(hdr sbus.MessageHandler) {
				if err := hdr(*msg, t.log); err != nil {
					t.log.WithError(err).Error("Error on handle message %v", msg)
				}
			}(handler)
		}
	}
	return nil
}
