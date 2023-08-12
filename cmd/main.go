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
	// зачем вам нужен URLFormat?
	// кажется его ошибочно использовал Тузов в своем гайде про url-shortener, видимо оттуда же)
	router.Use(middleware.URLFormat)

	containerAuth := BuildContainerAuthModule(connect, cfg)

	// неактуальная лог запись, тут еще ничего не стартануло
	log.Info("Success start")

	err := containerAuth.Invoke(func(handler *handler.RestHandler) {
		handler.HandleRequests(router)
	})
	if err != nil {
		log.Error("Rest handle crush : %v", err)
	}

	// тут форматинг не работает
	// будет такая строка: msg="Start on address: %s" !BADKEY=localhost:8001
	log.Info("Start on address: %s", cfg.HttpConf.Host+":"+cfg.HttpConf.Port)

	srv := &http.Server{
		// в конфиге можно сразу держать строку вида host:port
		Addr:         cfg.HttpConf.Host + ":" + cfg.HttpConf.Port,
		Handler:      router,
		ReadTimeout:  100 * time.Millisecond,
		WriteTimeout: 100 * time.Millisecond,
		IdleTimeout:  100 * time.Millisecond,
	}

	// нет graceful shutdown — аккуратного завершения работы по SIGINT/SIGTERM
	err = srv.ListenAndServe()
	if err != nil {
		log.Error("failed to listen: %v", err)
		// return?
	}

	// трейсер инициализируется уже после того как сервер завершил свою работу
	// получается трейсер не работает
	err = tracer.NewTracer("http://jaeger:14268/api/traces", "server")
	if err != nil {
		log.Error("error with tracer", err)
	}

	// тоже можно отнести к отсусттвию graceful shutdown, но тут следует задавать таймаут
	defer tracer.Tracer.Shutdown(context.Background())
}

// DI-контейнеры и gorm это интересно, но это все нужно чтобы уменьшить boilerplate в больших проектах
// их использование в go часто хейтится — они увеличивают learning curve и усложняют отладку
// так что вам стоит подумать, стоит ли это того
// особенно советую отказаться от GORM
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

	// вот тут, например, не провеяется ошибка
	// значит в проекте не настроены линтеры, даже базовые
	// Это первое что нужно сделать: https://golangci-lint.run
	// советую при обучении включать как можно больше линтеров
	err = container.Provide(handler.NewHandler)

	return container
}
