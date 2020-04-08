package hermes

import (
	"encoding/json"
	"errors"
	"github.com/go-pg/pg/v9"
	"io"
	"time"
)

var (
	db       *pg.DB
	messages chan Message
	quit     chan bool
	relays   []Relay
)

func RegisterRelay(relay Relay) {
	relays = append(relays, relay)
}

func InitStorage(dbName, dbHost, dbUser, dbPass string) {
	db = pg.Connect(&pg.Options{
		User:     dbUser,
		Password: dbPass,
		Database: dbName,
	})
}

func InitQueue() {
	messages = make(chan Message, 30)
	quit = make(chan bool)
	go MQ(messages, quit)
}

func CloseQueue() {
	quit <- true
}

func CloseStorage() {
	db.Close()
}

func MQ(msgs chan Message, ctl chan bool) {

}

type Customer struct {
	Id      int64
	Name    string
	Contact string
	Token   string `pg:"type:uuid,default:gen_random_uuid(),key"`
	// relations
	Services []*Service
}

type Service struct {
	CustomerId int64   `pg:",pk"`
	Type       [2]byte `pg:",pk"`
	Url        string
	Username   string
	Password   string
	Token      string
	Key        string
	Cert       string
	// relations
	Customer *Customer
}

type Relay interface {
	Code() [2]byte
	Send(message *Message, settings Service) error
	CanSendTo(recipient string) bool
}

func GetCustomer(token string) (Customer, error) {
	var (
		customer Customer
	)
	if token == "" {
		return customer, errors.New("Token missing")
	}
	err := db.Model(&customer).Where("token = ?", token).Relation("Services").Select()
	return customer, err
}

type Message struct {
	Uuid       string `pg:"type:uuid,default:gen_random_uuid(),pk"`
	Recipient  string
	Subject    string
	Message    string
	Failures   int
	Created    time.Time
	Sheduled   time.Time // when to send the message
	Sent       time.Time
	Status     string
	CustomerId int64
	Type       [2]byte
	// relations
	Customer *Customer
}

func ParseMessage(r io.Reader) (Message, error) {
	var (
		message Message
	)

	decoder := json.NewDecoder(r)
	err := decoder.Decode(&message)
	return message, err
}

func Send(customer Customer, message Message) (string, error) {
	err := db.Insert(&message)
	if err != nil {
		return "", err
	}
	messages <- message
	return message.Uuid, err
}

// type Reply struct {
// 	Message  string
// 	Sent     time.Time
// 	Received time.Time
// 	Sender   string
// }

// type Status struct {
// 	Uuid      string
// 	Scheduled time.Time
// 	Sent      time.Time
// 	Failures  int
// 	Replies   []Reply
// }

// func GetStatus(customer *Customer, uuid string) Status {

// }
