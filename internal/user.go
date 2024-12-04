package internal

type User struct {
	Login    string
	Password string
	Balance  int64
}

type Token struct {
	Login      string
	Password   string
	Expiration string
	Issued     string
}
