package auth

import (
	"auth-service/pkg/grpc/client"
	"auth-service/pkg/grpc/userGrpc"
	"context"
	"strconv"
)

type AuthService struct {
	repo           *DbRepository
	grpcUserClient *client.UserClientGrpc
}

func (s *AuthService) CreateAuth(auth Auth) {
	createAuth := s.repo.CreateAuth(auth)
	//TODO add logic to check exist user or create new
	s.grpcUserClient.Client.GetUserById(context.Background(), &userGrpc.GetUserById_Request{
		UserId: strconv.Itoa(int(createAuth.ID))})
}

func NewService(repo *DbRepository, grpcUserClient *client.UserClientGrpc) *AuthService {
	return &AuthService{repo: repo, grpcUserClient: grpcUserClient}
}
