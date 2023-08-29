package domain

import "context"

type Service interface {
	CreateAuth(ctx context.Context, auth AuthRepo) (createdAuth AuthFetch, err error)
	FetchAuth(ctx context.Context, id, email string) (AuthFetch, error)
	ResetPassword(ctx context.Context, id, email string) bool
	DeleteAuth(ctx context.Context, id string) bool
	UpdateAuth(ctx context.Context, auth AuthUpdate) (AuthUpdate, error)
}

type Repository interface {
	CreateAuth(ctx context.Context, auth AuthRepo) (createAuth AuthFetch, err error)
	DeleteAuth(ctx context.Context, id string) bool
	FetchAuth(ctx context.Context, id, email string) (auth AuthFetch, err error)
	UpdateAuth(ctx context.Context, id string, auth AuthUpdate) (AuthUpdate, error)
}
