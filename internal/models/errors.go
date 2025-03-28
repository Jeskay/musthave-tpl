// Module models describes objects used in business logic of gophermart loyalty system.
package models

import "errors"

// ErrUsedLogin indicates that login for registration was already used.
var ErrUsedLogin = errors.New("user with the same login already exists")

// ErrIncorrectPassword indicates that client used incorrect login or password.
var ErrIncorrectPassword = errors.New("incorrect login/password pair")

// ErrOrderExists indicates that order with current ID is already registered in the loyalty system.
var ErrOrderExists = errors.New("order with provided id already exists")

// ErrOrderUsed indicates that order with current ID was already used by another user.
var ErrOrderUsed = errors.New("order with provided id already used by other user")

// ErrInvalidOrderNumber indicates that current order ID is incorrect.
var ErrInvalidOrderNumber = errors.New("incorrect order id")

// ErrNotEnoughFunds  indicates that user balance does not have enough funds to process the withdrawal.
var ErrNotEnoughFunds = errors.New("not enough funds to process the request")
