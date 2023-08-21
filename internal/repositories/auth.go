package repositories

import (
	authDomain "auth-service/internal/auth"
	"context"
	"errors"
	"fmt"
	"golang.org/x/exp/slog"
	"gorm.io/gorm"
	"strconv"
)

type DbRepository struct {
	connect *gorm.DB
	logger  *slog.Logger
}

func (repo *DbRepository) CreateAuth(ctx context.Context, auth authDomain.AuthRepoStruct) (createAuth authDomain.AuthRepoStruct, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = authDomain.ErrorWhileCreate
		}
	}()
	//TODO constraint on email model
	//TODO add logic to hash pass
	res := repo.connect.Create(&auth).WithContext(ctx)

	if res.RowsAffected == 0 {
		return auth, errors.New("no update for this query")
	}
	return auth, nil
}

func (repo *DbRepository) DeleteAuth(ctx context.Context, id string) {
	repo.connect.Delete(&authDomain.AuthRepoStruct{}, id).WithContext(ctx)
}

func (repo *DbRepository) FetchAuth(ctx context.Context, id, email string) (auth authDomain.AuthRepoStruct, err error) {
	if id != "" {
		idUuid, err := strconv.ParseUint(id, 10, 32)
		if err != nil {
			return auth, fmt.Errorf("Error while transform id %s", err)
		}
		res := repo.connect.Find(auth, "id = ?", uint(idUuid)).WithContext(ctx)
		if res.Error != nil {
			repo.logger.Error(res.Error.Error())
			return authDomain.AuthRepoStruct{}, authDomain.ErrorNoAuth
		}
		return auth, nil
	}
	if email != "" {
		auth := authDomain.AuthRepoStruct{}
		res := repo.connect.Find(&auth, "email = ?", email)
		if res.Error != nil {
			repo.logger.Error(res.Error.Error())
			return authDomain.AuthRepoStruct{}, authDomain.ErrorNoAuth
		}
		return auth, nil
	}
	return authDomain.AuthRepoStruct{}, authDomain.ErrorNoAuth
}

func (repo *DbRepository) UpdateAuth(ctx context.Context, id string, auth authDomain.AuthRepoStruct) (authDomain.AuthRepoStruct, error) {
	res := repo.connect.Model(&auth).Where("id = ?", id).Updates(auth).WithContext(ctx)
	if res.RowsAffected == 0 {
		return auth, errors.New("no update for this query")
	}
	return auth, nil
}

func NewRepository(db *gorm.DB, logger *slog.Logger) authDomain.AuthRepository {
	return &DbRepository{connect: db, logger: logger}
}
