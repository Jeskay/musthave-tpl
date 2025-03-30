package models

// User - gophermart loyalty system user. Stores information about accrued and withdrawn bonuses.
type User struct {
	Login     string
	Password  string
	Balance   float64
	Withdrawn float64
}

// Token - user token. Holds expiration date and user information.
type Token struct {
	Login      string
	Password   string
	Expiration string
	Issued     string
}
