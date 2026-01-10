package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"pipe-scaling-analyzer/internal/config"
	"pipe-scaling-analyzer/internal/database"
	"pipe-scaling-analyzer/internal/httpserver"
	"pipe-scaling-analyzer/internal/logger"
	"pipe-scaling-analyzer/internal/provider"
	"pipe-scaling-analyzer/internal/repository"
	"pipe-scaling-analyzer/internal/usecase"
	"syscall"
	"time"
)

func main () {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	zapLog,closelog,err := logger.NewLogger(cfg.LogLevel)
	if err != nil {
		log.Fatalf("failed to create logger: %v", err)
	}
	defer closelog()
	
	ctx := context.Background()
	
	dbPool,err := database.NewPostgresPool(ctx, cfg.PostgresDSN)
	if err != nil {
		log.Fatalf("failed to connect postgres: %v",err)
	}
	defer dbPool.Close()



	pgRepo := repository.NewPostgresRepo(dbPool,zapLog)
	redisCache := repository.NewRedisCache(cfg.RedisAddr,cfg.RedisPass,zapLog)
	
	calcProv := provider.NewCalcModelProvider(cfg.CalcModelURL,zapLog)

	us := usecase.NewUseCase(pgRepo,redisCache,calcProv,zapLog)
	
	handlers := httpserver.NewHTTPHandlers(us,zapLog)
	srv := httpserver.NewServer(cfg.HTTPPort,handlers)

	go func() {
		zapLog.Infof("pipe-scaling-analyzer listening on :%s", cfg.HTTPPort)
		
		if err := srv.Start(); err != nil {
			log.Fatalf("failed to serve %v", err)
		}
	}()


	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	
	<-sigChan

	zapLog.Infof("shutting down server...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	_ = srv.Shutdown(shutdownCtx)
}
