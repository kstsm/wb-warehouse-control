package cmd

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kstsm/wb-warehouse-control/config"
	"github.com/kstsm/wb-warehouse-control/database"
	"github.com/kstsm/wb-warehouse-control/internal/handler"
	"github.com/kstsm/wb-warehouse-control/internal/repository"
	"github.com/kstsm/wb-warehouse-control/internal/service"
	"github.com/kstsm/wb-warehouse-control/pkg/jwt"
	"github.com/kstsm/wb-warehouse-control/pkg/logger"
	"github.com/kstsm/wb-warehouse-control/pkg/validator"
)

const (
	httpServerShutdownTimeout = 5
	readHeaderTimeout         = 5
)

func Run() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg := config.GetConfig()
	log := logger.NewSlogLogger()

	conn := database.InitPostgres(ctx, cfg, log)
	defer conn.Close()

	validate := validator.NewValidator()

	tokenManager := jwt.NewJWTManager(cfg.JWT.Secret, cfg.JWT.TTL, cfg.JWT.Issuer)

	repo := repository.NewRepository(conn, log)
	svc := service.NewService(repo, log, tokenManager)
	router := handler.NewHandler(svc, log, validate, tokenManager)

	srv := &http.Server{
		Addr:              fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler:           router.NewRouter(),
		ReadHeaderTimeout: readHeaderTimeout * time.Second,
	}

	errChan := make(chan error, 1)

	go func() {
		log.Infof("Starting server on %s:%d", cfg.Server.Host, cfg.Server.Port)
		errChan <- srv.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		log.Info("Finishing the server...")
	case err := <-errChan:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Errorf("Error starting server: %v", err)
			return err
		}
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), httpServerShutdownTimeout*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Errorf("Error while shutting down the server: %v", err)
	}

	return nil
}
