package relay

import (
	"bytes"
	"context"
	"gopkg.in/gomail.v2"
	"hermes/queue"
	"net/smtp"
	"net/url"
	"strings"
)

const (
	EMAIL_RELAY_CODE = "@"
)

var (
	barricade = make(chan bool, 3)
)

type EmailRelay Relay

func (relay EmailRelay) Send(ctx context.Context, message *queue.Message) error {
	var (
		done chan error
		err  error
	)
	select {
	case barricade <- true:
		done = relay.send(ctx, message)
	case <-ctx.Done():
		return ctx.Err()
	}
	defer func() { <-barricade }()

	select {
	case err = <-done:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (relay EmailRelay) CanSendTo(recipient string) bool {
	return true
}

func (relay EmailRelay) send(ctx context.Context, message *queue.Message) chan error {
	done := make(chan error)
	go func() {
		var (
			err error
			url *url.URL
		)
		defer func() { done <- err }()

		url, err = url.Parse(relay.Url)
		if err != nil {
			return
		}
		email := gomail.NewMessage()
		if message.Subject != "" {
			email.SetHeader("Subject", message.Subject)
		}
		if strings.Contains(message.Message, "<html>") {
			email.SetBody("text/html", message.Message)
		} else {
			email.SetBody("text/plain", message.Message)
		}
		buffer := new(bytes.Buffer)
		_, err = email.WriteTo(buffer)
		if err != nil {
			return
		}
		err = smtp.SendMail(
			relay.Url,
			smtp.PlainAuth("", relay.Username, relay.Password, url.Hostname()),
			relay.From,
			[]string{message.Recipient},
			buffer.Bytes(),
		)
	}()
	return done
}
