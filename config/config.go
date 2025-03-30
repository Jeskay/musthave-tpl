// Module config describes gophermart configuration.
package config

import (
	"crypto/rand"
	"time"
)

// Config - gophermart server configuration.
type Config struct {
	// Address - server address to run on.
	Address string `env:"RUN_ADDRESS"`
	// AccrualAddress - accrual service address.
	AccrualAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
	// HashKey - hash key for user passwords. Optional.
	HashKey string `env:"HASH_SECRET"`
	// TokenKey - JWT key for authentication. Optional.
	TokenKey string `env:"TOKEN_SECRET"`
	// TokenExpire - JWT expiration time. In milliseconds.
	TokenExpire int64 `env:"TOKEN_EXPIRE"`
	// DBConnection - database connection URI.
	DBConnection string `env:"DATABASE_URI"`
}

// NewGophermartConfig creates new Config instance.
func NewGophermartConfig() *Config {
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
	return &Config{
		TokenKey:    string(token1),
		TokenExpire: int64(time.Hour * 48),
		HashKey:     string(token2),
	}
}
