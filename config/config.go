package config

type LoyaltyConfig struct {
	Address      string `env:"ADDRESS"`
	HashKey      string `env:"HASH_SECRET"`
	TokenKey     string `env:"TOKEN_SECRET"`
	DBConnection string `env:"DATABASE_DSN"`
}

func NewLoyaltyConfig() *LoyaltyConfig {
	return &LoyaltyConfig{}
}
