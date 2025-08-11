package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"wallet/internal/config"
	"wallet/internal/handlers/wallet"
	"wallet/pkg/customValidator"
	"wallet/pkg/logger"
	"wallet/storage/postgres"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {

	logger := logger.InitLogger("./logs")
	logger.Info("logger init")

	config := config.New("./config.env")
	logger.Info("config init")

	storage := postgres.New(config)
	logger.Info("Successfully connected to postgres")

	cv := customValidator.NewCustomValidator()
	logger.Info("custom validator init")

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())

	router.POST("api/v1/wallet", wallet.ChangeAmountWallet(logger, cv, storage))
	router.GET("api/v1/wallet/:id", wallet.Amount(logger, cv, storage))

	srv := &http.Server{
		Addr:    ":" + config.AppPort,
		Handler: router.Handler(),
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	logger.Info("server start", zap.String("port", config.AppPort))
	gracefulShutdown(srv)
}

func gracefulShutdown(srv *http.Server) {

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Println("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}
