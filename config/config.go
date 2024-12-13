package config

type LoyaltyConfig struct {
	Address        string `env:"RUN_ADDRESS"`
	AccrualAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
	HashKey        string `env:"HASH_SECRET"`
	TokenKey       string `env:"TOKEN_SECRET"`
	DBConnection   string `env:"DATABASE_DSN"`
}

func NewLoyaltyConfig() *LoyaltyConfig {
	return &LoyaltyConfig{}
}
