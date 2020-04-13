package core

import (
	"errors"
	"hermes/storage"
)

type Customer struct {
	Id      int64 `pg:",pk"`
	Name    string
	Contact string
	Token   string `pg:"type:uuid,default:gen_random_uuid(),key"`
}

func GetCustomer(token string) (Customer, error) {
	var (
		customer Customer
	)
	if token == "" {
		return customer, errors.New("Token missing")
	}
	err := storage.Get(&customer, "token = ?", token)
	return customer, err
}
