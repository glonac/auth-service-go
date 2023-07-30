package client

import (
	"auth-service/internal/config"
	"auth-service/internal/grpc/userGrpc"
	"auth-service/internal/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type UserClientGrpc interface {
	userGrpc.UserServiceClient
}

func NewUserClient(cnf *config.ConfigGrpc) UserClientGrpc {
	conn, err := grpc.Dial(cnf.Host+":"+cnf.Port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Err(err)
	}

	client := userGrpc.NewUserServiceClient(conn)
	return client
}
