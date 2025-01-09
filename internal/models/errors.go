package models

type UsedLoginError struct{}

func (usedLogin *UsedLoginError) Error() string {
	return "user with the same login already exists"
}

type IncorrectPassword struct{}

func (incorrectPass *IncorrectPassword) Error() string {
	return "incorrect login/password pair"
}

type OrderExists struct{}

func (orderExist *OrderExists) Error() string {
	return "order with provided id already exists"
}

type OrderUsed struct{}

func (orderUsed *OrderUsed) Error() string {
	return "order with provided id already used by other user"
}

type InvalidOrderNumber struct{}

func (invalidNumber *InvalidOrderNumber) Error() string {
	return "incorrect order id"
}

type NotEnoughFunds struct{}

func (noFunds *NotEnoughFunds) Error() string {
	return "not enough funds to process the request"
}
