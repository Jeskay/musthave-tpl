package dto

// User - user, registered in accrual system. Describes JSON representation.
type User struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// Balance - accrual balance. Describes JSON representation.
type Balance struct {
	Balance   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}
