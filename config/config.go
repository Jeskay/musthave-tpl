package config

import "crypto/rand"

type GophermartConfig struct {
	Address        string `env:"RUN_ADDRESS"`
	AccrualAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
	HashKey        string `env:"HASH_SECRET"`
	TokenKey       string `env:"TOKEN_SECRET"`
	DBConnection   string `env:"DATABASE_URI"`
}

func NewGophermartConfig() *GophermartConfig {
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
	return &GophermartConfig{
		TokenKey: string(token1),
		HashKey:  string(token2),
	}
}
