package gophermart

type UsedLoginError struct{}

func (usedLogin *UsedLoginError) Error() string {
	return "user with the same login already exists"
}

type IncorrectPassword struct{}

func (incorrectPass *IncorrectPassword) Error() string {
	return "incorrect login/password pair"
}
