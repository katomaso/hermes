package relay

import (
	"context"
	"errors"
	"hermes/core"
	"hermes/queue"
	"hermes/storage"
	"log"
	"time"
)

type Relay struct {
	CustomerId int64  `pg:",pk"`
	Code       string `pg:"type:char(2),pk"`
	From       string
	Url        string
	Username   string
	Password   string
	Token      string
	Key        string
	Cert       string
	// relations
	Customer *core.Customer
}

type Relayer interface {
	Send(ctx context.Context, message *queue.Message) error
	CanSendTo(recipient string) bool
}

var (
	NoRelayFound = errors.New("No Relay Found")
	background   = context.Background()
	quit         = make(chan bool)
)

func GetRelayer(message queue.Message) (Relayer, error) {
	var (
		relays []*Relay
	)
	err := storage.Get(&relays, "relay.customer_id = ?", message.CustomerId)
	if err != nil {
		return nil, err
	}
	if len(relays) == 0 {
		return nil, NoRelayFound
	}
	for _, relay := range relays {
		if relay.Code == EMAIL_RELAY_CODE {
			email := EmailRelay(*relay)
			if email.CanSendTo(message.Recipient) {
				return email, nil
			}
		}
	}
	return nil, NoRelayFound
}

func InitRelay(mq chan queue.Message) {
	log.Println("Registering Relays to storage")
	storage.CreateStorage((*Relay)(nil))
	log.Println("Starting relay clients")
	quit = make(chan bool)
	go relay(mq)
}

func CloseRelay() {
	quit <- true
}

func relay(mq chan queue.Message) {
	for {
		select {
		case msg := <-mq:
			go relayMessage(msg)
		case <-quit:
			break
		}
	}
}

func relayMessage(msg queue.Message) {
	relayer, err := GetRelayer(msg)
	if err != nil {
		msg.SetError(err)
		return
	}
	ctx, done := context.WithTimeout(background, 15*time.Second)
	defer done()

	err = relayer.Send(ctx, &msg)
	if err != nil {
		msg.SetError(err)
		return
	}
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
