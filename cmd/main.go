package main

import (
	auth "auth-service/pkg/auth"
	"auth-service/pkg/config"
	"auth-service/pkg/grpc/client"
	"auth-service/pkg/handler"
	"auth-service/pkg/logger"
	mwLogger "auth-service/pkg/middleware/http"
	"auth-service/pkg/storage"
	"auth-service/pkg/tracer"
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/dig"
	"gorm.io/gorm"
	"net/http"
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
		Addr:    cfg.HttpConf.Host + ":" + cfg.HttpConf.Port,
		Handler: router,
		//ReadTimeout:  1,
		//WriteTimeout: 1,
		//IdleTimeout:  1,
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
	err := container.Provide(func() *auth.DbRepository {
		return auth.NewRepository(connectDb)
	})

	if err != nil {
		fmt.Println(err)
	}

	err = container.Provide(func() *client.UserClientGrpc {
		return client.NewUserClient(cnf.GrpcConf)
	})

	if err != nil {
		fmt.Println(err)
	}

	err = container.Provide(auth.NewService)

	if err != nil {
		fmt.Println(err)
	}

	err = container.Provide(handler.NewHandler)

	return container
}
