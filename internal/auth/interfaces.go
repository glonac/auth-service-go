package auth

import "context"

type AuthService interface {
	CreateAuth(ctx context.Context, auth AuthRepoStruct) (createdAuth AuthRepoStruct, err error)
	FetchAuth(ctx context.Context, id, email string) (AuthRepoStruct, error)
	ResetPassword(ctx context.Context, id, email string) bool
}

type AuthRepository interface {
	CreateAuth(ctx context.Context, auth AuthRepoStruct) (createAuth AuthRepoStruct, err error)
	DeleteAuth(ctx context.Context, id string)
	FetchAuth(ctx context.Context, id, email string) (auth AuthRepoStruct, err error)
	UpdateAuth(ctx context.Context, id string, auth AuthRepoStruct) (AuthRepoStruct, error)
}
