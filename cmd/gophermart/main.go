package main

import (
	"context"
	"database/sql"
	"flag"
	"log"
	"musthave_tpl/config"
	"musthave_tpl/internal"
	"musthave_tpl/internal/auth"
	"musthave_tpl/internal/gophermart"
	"musthave_tpl/internal/gophermart/db"
	"musthave_tpl/internal/loyalty"
	"musthave_tpl/internal/middleware"
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

	zapL, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}
	defer zapL.Sync()
	conf, err := envParse()
	if err != nil {
		zapL.Fatal("failed to parse environment variables", zap.Error(err))
	}
	server := initHelper(conf, zapL)
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zapL.Fatal("listen: ", zap.String("", err.Error()))
		}
	}()
	shutdownHelper(context.Background(), zapL, server)
}

func initHelper(conf config.Config, logger *zap.Logger) *http.Server {
	database, err := sql.Open("pgx", conf.DBConnection)
	if err != nil {
		logger.Fatal("failed to connect to the database", zap.Error(err))
	}
	repository, err := db.NewPostgresStorage(database, zapslog.NewHandler(logger.Core()))
	if err != nil {
		logger.Fatal("failed to initialize database", zap.Error(err))
	}
	authService := auth.NewAuthService(&conf)
	loyltyService := loyalty.NewLoyaltyService(&conf, zapslog.NewHandler(logger.Core()))
	middlewareService := middleware.NewMiddlewareService(authService, repository)
	gophermartService := gophermart.NewGophermartService(
		&conf,
		zapslog.NewHandler(logger.Core()),
		repository,
		authService,
		loyltyService,
	)
	app := internal.NewApp(
		gophermartService,
		middlewareService,
	)
	r := app.Router()
	return &http.Server{
		Addr:    conf.Address,
		Handler: r.Handler(),
	}
}

func shutdownHelper(ctx context.Context, logger *zap.Logger, server *http.Server) {
	sig := make(chan os.Signal, 1)

	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	logger.Info("shutdown server ...")

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal("server shutdown", zap.String("", err.Error()))
	}
	<-ctx.Done()

	logger.Info("server exiting")
}

func envParse() (config.Config, error) {
	conf := config.NewGophermartConfig()
	flag.StringVar(&conf.Address, "a", "", "service run address")
	flag.StringVar(&conf.DBConnection, "d", "", "database connection address")
	flag.StringVar(&conf.AccrualAddress, "r", "", "accrual service address")
	flag.Int64Var(&conf.TokenExpire, "e", int64(time.Hour*48), "token expiration time")
	flag.Parse()

	err := env.Parse(&conf)
	return conf, err
}
