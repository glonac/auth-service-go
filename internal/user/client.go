package user

import (
	"auth-service/internal/grpc/userGrpc"
	"context"
	"google.golang.org/grpc"
)

type UserClient interface {
	GetUserByClientAccountId(ctx context.Context, in *userGrpc.GetUserByClientAccountId_Request, opts ...grpc.CallOption) (*userGrpc.GetUserByClientAccountId_Response, error)
}
