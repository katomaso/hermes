package queue

import (
	"context"
	"hermes/core"
	"hermes/storage"
	"log"
)

var (
	messages chan Message
	quit     chan bool
)

func InitQueue() {
	log.Println("Registering Message to storage")
	storage.CreateStorage((*Message)(nil))
	log.Println("Starting a message queue")
	messages = make(chan Message, 30)
	quit = make(chan bool)
}

func CloseQueue() {
	quit <- true
}

func Queue() chan<- Message {
	return messages
}

/** Send accepts `message` and returns its identifier that can be later on used
 * for status queries.
 */
func Send(ctx *context.Context, customer *core.Customer, message Message) (string, error) {
	err := storage.Insert(&message)
	if err != nil {
		return "", err
	}
	messages <- message
	return message.Uuid, err
}

// func GetStatus(customer *Customer, uuid string) Status {

// }
