package main

import (
	auth "auth-service/internal/auth"
	"auth-service/internal/config"
	"auth-service/internal/grpc/client"
	"auth-service/internal/handler"
	"auth-service/internal/logger"
	mwLogger "auth-service/internal/middleware/http"
	"auth-service/internal/storage"
	"auth-service/internal/tracer"
	"auth-service/internal/user"
	"context"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/dig"
	"gorm.io/gorm"
)

func main() {
	cfg := config.MustLoad()
	log := logger.SetupLogger("local")
	connect := storage.Initialize(cfg.DataBaseConf)

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(mwLogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	containerAuth := BuildContainerAuthModule(connect, cfg)

	log.Info("Success start")

	err := containerAuth.Invoke(func(handler *handler.RestHandler) {
		handler.HandleRequests(router)
	})
	if err != nil {
		log.Error("Rest handle crush : %v", err)
	}

	log.Info("Start on address: %s", cfg.HttpConf.Host+":"+cfg.HttpConf.Port)

	srv := &http.Server{
		Addr:         cfg.HttpConf.Host + ":" + cfg.HttpConf.Port,
		Handler:      router,
		ReadTimeout:  100 * time.Millisecond,
		WriteTimeout: 100 * time.Millisecond,
		IdleTimeout:  100 * time.Millisecond,
	}
	err = srv.ListenAndServe()
	if err != nil {
		log.Error("failed to listen: %v", err)
	}

	err = tracer.NewTracer("http://jaeger:14268/api/traces", "server")
	if err != nil {
		log.Error("error with tracer", err)
	}

	defer tracer.Tracer.Shutdown(context.Background())
}

func BuildContainerAuthModule(connectDb *gorm.DB, cnf *config.MainConfig) *dig.Container {
	container := dig.New()
	err := container.Provide(func() auth.AuthRepository {
		return auth.NewRepository(connectDb)
	})

	if err != nil {
		log.Println(err)
	}

	err = container.Provide(func() user.UserClient {
		return client.NewUserClient(cnf.GrpcConf)
	})

	if err != nil {
		log.Println(err)
	}

	err = container.Provide(auth.NewService)

	if err != nil {
		log.Println(err)
	}

	err = container.Provide(handler.NewHandler)

	return container
}
