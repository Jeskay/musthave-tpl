package dto

type User struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Balance struct {
	Balance   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}
