package main

import (
	"context"
	"database/sql"
	"flag"
	"log"
	"musthave_tpl/config"
	"musthave_tpl/internal/gophermart"
	"musthave_tpl/internal/gophermart/db"
	"musthave_tpl/internal/gophermart/routes"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/caarlos0/env"
	"go.uber.org/zap"
	"go.uber.org/zap/exp/zapslog"
)

func main() {
	var conf config.GophermartConfig
	zapL, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}
	defer zapL.Sync()

	flag.StringVar(&conf.Address, "a", "", "service run address")
	flag.StringVar(&conf.DBConnection, "d", "", "database connection address")
	flag.StringVar(&conf.AccrualAddress, "r", "", "accrual service address")
	flag.Int64Var(&conf.TokenExpire, "e", int64(time.Hour*48), "token expiration time")
	flag.Parse()

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
	server := &http.Server{
		Addr:    conf.Address,
		Handler: r.Handler(),
	}
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zapL.Fatal("listen: ", zap.String("", err.Error()))
		}
	}()
	sig := make(chan os.Signal, 1)

	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	zapL.Info("shutdown server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		zapL.Fatal("server shutdown", zap.String("", err.Error()))
	}
	<-ctx.Done()

	zapL.Info("server exiting")
}
