package sbus

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/satori/go.uuid"
)

func New(transp Transport, log *logrus.Entry) *Sbus {
	return &Sbus{log, transp}
}

type Message struct {
	Subject string          `json:"subject"`
	Data    json.RawMessage `json:"data,omitempty"`
	Meta    Meta            `json:"meta,omitempty"`
}

func (m Message) WithMeta(key string, value interface{}) Message {
	if m.Meta == nil {
		m.Meta = Meta{}
	}
	m.Meta[key] = value
	return m
}

func (m Message) Unmarshal(v interface{}) error {
	return json.Unmarshal([]byte(m.Data), v)
}

func (m Message) String() string {
	re, _ := json.Marshal(m)
	return string(re)
}

type Meta map[string]interface{}

type MessageHandler func(Message) error

type Transport interface {
	Sub(subject string, handler MessageHandler) error
	SubOnce(subject string, handler MessageHandler) error
	Pub(msg *Message) error
}

type Sbus struct {
	log    *logrus.Entry
	transp Transport
}

func (s *Sbus) Sub(subject string, handler MessageHandler) {
	if err := s.transp.Sub(subject, handler); err != nil {
		s.log.WithError(err).Errorf("Error on subscribe to %s", subject)
	}
}

func (s *Sbus) Pub(subject string, data interface{}) error {
	return s.PubM(Message{subject, Marshal(data), nil})
}

func (s *Sbus) PubM(msg Message) error {
	err := s.transp.Pub(&msg)
	if err != nil {
		s.log.WithError(err).Errorf("Error on publish %v", msg)
	}
	return err
}

func (s *Sbus) Request(subject string, data interface{}, handler MessageHandler, timeout time.Duration) error {
	return s.RequestM(Message{subject, Marshal(data), nil}, handler, timeout)
}

func (s *Sbus) RequestM(msg Message, handler MessageHandler, timeout time.Duration) error {
	uid, err := uuid.NewV4()
	if err != nil {
		return err
	}
	
	replyTo := msg.Subject + "-" + uid.String()

	s.transp.SubOnce(replyTo, handler)
	return s.PubM(msg.WithMeta("replyTo", replyTo))
}

func (s *Sbus) Reply(msg Message, response interface{}) error {
	if replyTo, ok := msg.Meta["replyTo"]; ok {
		return s.Pub(fmt.Sprintf("%s", replyTo), response)
	}
	return fmt.Errorf("Error on replay: not found 'replyTo' field in request %v!", msg)
}

func Marshal(data interface{}) []byte {
	resp, _ := json.MarshalIndent(data, "", "  ")
	return resp
}
