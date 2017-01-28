package sbus

import (
	"fmt"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/satori/go.uuid"
)

func New(transp Transport, log *logrus.Entry) *Sbus {
	return &Sbus{log, transp}
}

type Message struct {
	Subject string
	Data    interface{}
	Meta    Meta
}

type Meta map[string]interface{}

type MessageHandler func(Message, *logrus.Entry) error

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
		s.log.Errorf("Error on subscribe to %s", subject)
	}
}

func (s *Sbus) Pub(subject string, data interface{}) error {
	return s.PubM(Message{subject, data, nil})
}

func (s *Sbus) PubM(msg Message) error {
	err := s.transp.Pub(&msg)
	if err != nil {
		s.log.Errorf("Error on publish %v", msg)
	}
	return err
}

func (s *Sbus) Request(subject string, data interface{}, handler MessageHandler, timeout time.Duration) error {
	return s.RequestM(Message{subject, data, nil}, handler, timeout)
}

func (s *Sbus) RequestM(msg Message, handler MessageHandler, timeout time.Duration) error {
	replyTo := msg.Subject + "-" + uuid.NewV4().String()

	s.transp.SubOnce(replyTo, handler)
	msg.Meta["replyTo"] = replyTo
	return s.PubM(msg)
}

func (s *Sbus) Reply(msg Message, response interface{}) error {
	if replyTo, ok := msg.Meta["replyTo"]; ok {
		return s.Pub(fmt.Sprintf("%s", replyTo), response)
	}
	return fmt.Errorf("Error on replay: not found 'replyTo' field in request %v!", msg)
}
