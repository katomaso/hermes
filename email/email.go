package email

import (
	"hermes"
	"net/smtp"
)

const (
	EMAIL_TYPE = [2]byte{'@', 'M'}
)

type EmailRelay struct{}

func (relay *EmailRelay) Code() [2]byte {
	return EMAIL_TYPE
}

func (relay *EmailRelay) Send(message *Message, settings Service) error {
	return nil
}

func (relay *EmailRelay) CanSendTo(recipient string) bool {
	return true
}

func init() {
	hermes.RegisterRelay(EmailRelay{})
}
