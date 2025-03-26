package models

import "errors"

var ErrUsedLogin = errors.New("user with the same login already exists")
var ErrIncorrectPassword = errors.New("incorrect login/password pair")
var ErrOrderExists = errors.New("order with provided id already exists")
var ErrOrderUsed = errors.New("order with provided id already used by other user")
var ErrInvalidOrderNumber = errors.New("incorrect order id")
var ErrNotEnoughFunds = errors.New("not enough funds to process the request")
