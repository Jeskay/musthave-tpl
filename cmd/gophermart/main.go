package main

import (
	"database/sql"
	"musthave_tpl/config"
	"musthave_tpl/internal/gophermart"
	"musthave_tpl/internal/gophermart/db"
	"musthave_tpl/internal/gophermart/routes"

	"github.com/caarlos0/env"
	"go.uber.org/zap"
	"go.uber.org/zap/exp/zapslog"
)

func main() {
	var conf config.LoyaltyConfig
	zapL := zap.Must(zap.NewProduction())

	if err := env.Parse(&conf); err != nil {
		zapL.Fatal("failed to parse environment variables", zap.Error(err))
	}

	database, err := sql.Open("pgx", conf.DBConnection)
	if err != nil {
		zapL.Fatal("failed to connect to the database", zap.Error(err))
	}
	storage, err := db.NewPostgresStorage(database, zapslog.NewHandler(zapL.Core()))
	if err != nil {
		zapL.Fatal("failed to initialize database", zap.Error(err))
	}
	service := gophermart.NewGophermartService(&conf, zapslog.NewHandler(zapL.Core()), storage)

	r := routes.Init(service)

	r.Run(conf.Address)
}
