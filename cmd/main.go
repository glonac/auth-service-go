package main

import (
	auth "auth-service/pkg/auth"
	"auth-service/pkg/config"
	"auth-service/pkg/handler"
	"auth-service/pkg/logger"
	"auth-service/pkg/storage"
	"fmt"
	"net/http"
	"os"
	"gorm.io/gorm"
)

func main() {
  cfg := config.MustLoad()
  log := logger.SetupLogger("local")
	connect := storage.Initialize(cfg)
	migrations(connect)

	containerAuth := BuildContainerAuthModule(connect)

  log.Info("Success start")

	err := containerAuth.Invoke(func(handler *handler.RestHandler) {
		handler.HandleRequests()
	})
	if err != nil {
    log.Error("Rest hadler crush : %v", err)
	}

	port := ":" + os.Getenv("PORT")
	if os.Getenv("PORT") == "" {
		port = ":8080"
	}
  log.Info("Start on port: %s", port)

	err = http.ListenAndServe(port, nil)
	if err != nil {
		log.Error("failed to listen: %v", err)
	}
	//err := tracer.NewTracer("http://jaeger:14268/api/traces", "server")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//defer tracer.Tracer.Shutdown(context.Background())
}


func BuildContainerAuthModule(connectDb *gorm.DB) *dig.Container {
	container := dig.New()
	err := container.Provide(func() *auth.DbRepository {
		return auth.NewRepository(connectDb)
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

func migrations(connect *gorm.DB) {
	err := connect.AutoMigrate(&auth.Auth{})
	if err != nil {
		panic(err)
	}
}
