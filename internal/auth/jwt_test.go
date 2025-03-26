package auth

import (
	"math/rand"
	"musthave_tpl/config"
	"strconv"
	"testing"
)

func BenchmarkShortTokens(b *testing.B) {
	conf := config.NewGophermartConfig()
	service := NewAuthService(conf)
	count := 10000
	input := make([]string, count)
	for i := 0; i < count; i++ {
		input[i] = "string " + strconv.Itoa(i)
	}
	tokens := make([]string, count)
	var err error
	b.Run("creation", func(b *testing.B) {
		for i, l := range input {
			tokens[i], err = service.CreateToken(l)
			if err != nil {
				b.FailNow()
			}
		}
	})
	b.Run("verification", func(b *testing.B) {
		for i, t := range tokens {
			res, err := service.VerifyToken(t)
			if err != nil || res != input[i] {
				b.FailNow()
			}
		}
	})
}

func BenchmarkLongTokens(b *testing.B) {

	conf := config.NewGophermartConfig()
	service := NewAuthService(conf)
	count := 10000
	input := make([]string, count)
	for i := 0; i < count; i++ {
		input[i] = randStringBytes(1000)
	}
	tokens := make([]string, count)
	var err error
	b.Run("creation", func(b *testing.B) {
		for i, l := range input {
			tokens[i], err = service.CreateToken(l)
			if err != nil {
				b.FailNow()
			}
		}
	})
	b.Run("verification", func(b *testing.B) {
		for i, t := range tokens {
			res, err := service.VerifyToken(t)
			if err != nil || res != input[i] {
				b.FailNow()
			}
		}
	})
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
