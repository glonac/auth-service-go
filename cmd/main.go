package main

import (
	"auth-service/internal/auth"
	"auth-service/internal/closer"
	"auth-service/internal/config"
	"auth-service/internal/grpc/client"
	"auth-service/internal/handler/rest"
	"auth-service/internal/logger"
	mwLogger "auth-service/internal/middleware/http"
	"auth-service/internal/queue"
	"auth-service/internal/repositories"
	"auth-service/internal/storage"
	"auth-service/internal/tracer"
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/dig"
	"golang.org/x/exp/slog"
	"gorm.io/gorm"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

const shutdownTimeout = 5 * time.Second

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	cfg := config.MustLoad()
	log := logger.SetupLogger("local")
	connect := storage.Initialize(cfg.DataBaseConf)
	c := closer.GetInstance()

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(mwLogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	containerAuth := BuildContainerAuthModule(connect, log, cfg)

	log.Info("Success start")

	err := containerAuth.Invoke(func(handler *rest.Handler) {
		handler.HandleRequests(router)
	})

	if err != nil {
		log.Error("Rest handle crush : %v", err)
	}

	log.Info("Start on address: %s\n", cfg.HttpConf.Host+":"+cfg.HttpConf.Port)

	srv := &http.Server{
		Addr:         cfg.HttpConf.Host + ":" + cfg.HttpConf.Port,
		Handler:      router,
		ReadTimeout:  100 * time.Millisecond,
		WriteTimeout: 100 * time.Millisecond,
		IdleTimeout:  100 * time.Millisecond,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("listen and serve: %v", err)
		}
	}()

	c.Add(srv.Shutdown)

	err = containerAuth.Invoke(func(handler *rest.Handler) {
		handler.HandleRequests(router)
	})

	if err != nil {
		log.Error("error request handle", err)
	}
	err = tracer.NewTracer("http://jaeger:14268/api/traces", "server")

	if err != nil {
		log.Error("error with tracer", err)
	}

	defer tracer.Tracer.Shutdown(context.Background())
	if err != nil {
		log.Error("failed to listen: %v", err)
	}
	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	err = c.Close(shutdownCtx)
	if err != nil {
		log.Error("error while shutdown", err)
	}
}

func BuildContainerAuthModule(connectDb *gorm.DB, logger *slog.Logger, cnf *config.MainConfig) *dig.Container {
	container := dig.New()
	err := container.Provide(func() auth.AuthRepository {
		return repositories.NewRepository(connectDb, logger)
	})

	if err != nil {
		logger.Error("error while build", err)
	}

	err = container.Provide(func() *client.UserClientGrpc {
		return client.NewUserClient(cnf.GrpcConf)
	})

	if err != nil {
		logger.Error("error while build", err)
	}

	err = container.Provide(auth.NewService)

	if err != nil {
		logger.Error("error while build", err)
	}

	err = container.Provide(func() queue.QueueService {
		return queue.NewClientQueue(cnf.QueueConf)
	})

	if err != nil {
		logger.Error("error while build", err)
	}
	err = container.Provide(rest.NewHandler)

	if err != nil {
		logger.Error("error while build", err)
	}

	err = container.Provide(func() *slog.Logger {
		return logger
	})

	if err != nil {
		logger.Error("error while build", err)
	}

	return container
}
