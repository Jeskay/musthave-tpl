package models

type User struct {
	Login     string
	Password  string
	Balance   float64
	Withdrawn float64
}

type Token struct {
	Login      string
	Password   string
	Expiration string
	Issued     string
}
