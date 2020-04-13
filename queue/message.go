package queue

import (
	"encoding/json"
	"hermes/core"
	"hermes/storage"
	"io"
	"time"
)

type Message struct {
	Uuid       string `pg:"type:uuid,default:gen_random_uuid(),pk",json:"-"`
	Recipient  string
	Subject    string
	Message    string
	Failures   int       `json:"-"`
	Created    time.Time `json:"-"`
	Sheduled   time.Time `json:"-"` // when to send the message
	Sent       time.Time `json:"-"`
	Status     string    `json:"-"`
	CustomerId int64     `json:"-"`
	// relations
	Customer *core.Customer
}

func ParseMessage(r io.Reader) (Message, error) {
	var (
		message Message
	)
	decoder := json.NewDecoder(r)
	err := decoder.Decode(&message)
	return message, err
}

func (msg *Message) SetError(err error) {
	msg.Failures++
	msg.Status = "ERROR: " + err.Error()
	msg.Save()
}

func (msg *Message) Save() {
	if msg.Uuid == "" {
		storage.Insert(msg)
	} else {
		storage.Update(msg)
	}
}
