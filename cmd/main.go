package main

import (
	"auth-service/internal/closer"
	"auth-service/internal/config"
	"auth-service/internal/domain"
	"auth-service/internal/grpc/client"
	"auth-service/internal/handler/rest"
	routerCustom "auth-service/internal/handler/rest/router"
	"auth-service/internal/logger"
	mwLogger "auth-service/internal/middleware/http"
	"auth-service/internal/repositories"
	"auth-service/internal/storage"
	"auth-service/internal/tracer"
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"golang.org/x/exp/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const shutdownTimeout = 5 * time.Second

func main() {
	log := logger.SetupLogger("local")
	if err := run(log); err != nil {
		log.Error("Fatal:", err)
	}
	os.Exit(0)
}

func run(log *slog.Logger) error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg := config.MustLoad()
	connect, err := storage.NewPG(context.Background(), cfg.DataBaseConf)
	if err != nil {
		log.Error(err.Error())
	}
	c := closer.GetInstance()

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(mwLogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	routerRest := buildRestTransport(connect, log, cfg)

	log.Info("Success start")

	routerRest.HandleRequests(router)

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

	err = tracer.NewTracer("http://jaeger:14268/api/traces", "server")

	if err != nil {
		log.Error("error with tracer", err)
		return err
	}

	defer func(Tracer *tracesdk.TracerProvider, ctx context.Context) {
		err := Tracer.Shutdown(ctx)
		if err != nil {
			log.Error("main : ", err.Error())
		}
	}(tracer.Tracer, context.Background())

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	err = c.Close(shutdownCtx)
	if err != nil {
		log.Error("error while shutdown", err)
		return err
	}
	return nil
}

func buildRestTransport(connectDb *storage.Postgres, logger *slog.Logger, cnf *config.MainConfig) routerCustom.Router {
	authRepo := repositories.NewRepository(connectDb, logger)
	userClient := client.NewUserClient(cnf.GrpcConf)
	authService := domain.NewService(authRepo, userClient, logger)
	restHandler := rest.NewHandler(authService, logger)
	router := routerCustom.NewRouter(restHandler)
	return router
}
