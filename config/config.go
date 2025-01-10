package config

import "crypto/rand"

type Config struct {
	Address        string `env:"RUN_ADDRESS"`
	AccrualAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
	HashKey        string `env:"HASH_SECRET"`
	TokenKey       string `env:"TOKEN_SECRET"`
	TokenExpire    int64  `env:"TOKEN_EXPIRE"`
	DBConnection   string `env:"DATABASE_URI"`
}

func NewGophermartConfig() Config {
	token1 := make([]byte, 20)
	token2 := make([]byte, 20)
	_, err := rand.Read(token1)
	if err != nil {
		token1 = []byte("defaultToken")
	}
	_, err = rand.Read(token2)
	if err != nil {
		token2 = []byte("defaultHash")
	}
	return Config{
		TokenKey: string(token1),
		HashKey:  string(token2),
	}
}
