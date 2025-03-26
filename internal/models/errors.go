package models

import "errors"

var UsedLoginError = errors.New("user with the same login already exists")
var IncorrectPassword = errors.New("incorrect login/password pair")
var OrderExists = errors.New("order with provided id already exists")
var OrderUsed = errors.New("order with provided id already used by other user")
var InvalidOrderNumber = errors.New("incorrect order id")
var NotEnoughFunds = errors.New("not enough funds to process the request")
