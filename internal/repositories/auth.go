package repositories

import (
	"auth-service/internal/domain"
	"auth-service/internal/storage"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"golang.org/x/exp/slog"
	"strconv"
	"time"
)

type DbRepository struct {
	connect *storage.Postgres
	logger  *slog.Logger
}

func (repo *DbRepository) CreateAuth(ctx context.Context, auth domain.AuthRepo) (createAuth domain.AuthFetch, err error) {
	var authCreated domain.AuthFetch
	//TODO constraint on email model
	//TODO add logic to hash pass
	query := `INSERT INTO auths (email, code , password , code_created_at , is_verified, created_at , updated_at) 
				VALUES (@email, @code, @password, @codeCreatedAt, @isVerified, @createdAt,@updatedAt) 
				returning (id,email, code , code_created_at , is_verified, created_at , updated_at)`

	args := pgx.NamedArgs{
		"email":         auth.Email,
		"password":      auth.Password,
		"code":          auth.Code,
		"codeCreatedAt": auth.CodeCreatedAt,
		"isVerified":    auth.IsVerified,
		"createdAt":     auth.CreatedAt,
		"updatedAt":     auth.UpdatedAt,
	}

	rows := repo.connect.Db.QueryRow(ctx, query, args)
	err = rows.Scan(&authCreated)
	if err != nil {
		repo.logger.Error("auth repository:", err)
		return domain.AuthFetch{}, err
	}

	return authCreated, nil
}

func (repo *DbRepository) DeleteAuth(ctx context.Context, id string) bool {
	query := `DELETE FROM auths WHERE id = $1`
	_, err := repo.connect.Db.Exec(ctx, query, id)
	if err != nil {
		repo.logger.Error("auth repo:", err.Error())
		return false
	}
	return true
}

func (repo *DbRepository) FetchAuth(ctx context.Context, id, email string) (auth domain.AuthFetch, err error) {
	var authCreated domain.AuthFetch
	if id != "" {
		idUuid, err := strconv.ParseUint(id, 10, 32)
		if err != nil {
			return auth, fmt.Errorf("error while transform id %s", err)
		}

		query := `SELECT 
     		id, email, code , code_created_at , is_verified, created_at , updated_at
    		FROM auths WHERE id = $1 LIMIT 1`

		rows := repo.connect.Db.QueryRow(ctx, query, idUuid)
		err = rows.Scan(
			&authCreated.ID,
			&authCreated.Email,
			&authCreated.Code,
			&authCreated.CodeCreatedAt,
			&authCreated.IsVerified,
			&authCreated.CreatedAt,
			&authCreated.UpdatedAt,
		)

		if err != nil {
			repo.logger.Error("auth repository:", err)
			return auth, err
		}

		return authCreated, nil
	}
	if email != "" {
		auth := domain.AuthFetch{}

		query :=
			`SELECT 
     		id ,email, code , code_created_at , is_verified, created_at , updated_at
    	FROM auths WHERE email = $1 LIMIT 1`

		rows := repo.connect.Db.QueryRow(ctx, query, email)
		err = rows.Scan(
			&authCreated.ID,
			&authCreated.Email,
			&authCreated.Code,
			&authCreated.CodeCreatedAt,
			&authCreated.IsVerified,
			&authCreated.CreatedAt,
			&authCreated.UpdatedAt,
		)

		if err != nil {
			repo.logger.Error("auth repository:", err)
			return auth, err
		}

		return authCreated, nil
	}
	return domain.AuthFetch{}, domain.ErrorNoAuth
}

func (repo *DbRepository) UpdateAuth(ctx context.Context, id string, auth domain.AuthUpdate) (domain.AuthUpdate, error) {
	query := `UPDATE auths SET `
	if auth.Email != "" {
		query = query + `email = @email`
	}
	if auth.Code != "" {
		query = query + `, code = @code`
	}
	if auth.CodeCreatedAt != time.Now() {
		query = query + `, code_created_at= @codeCreatedAt`
	}
	if auth.IsVerified {
		query = query + `, is_verified= @isVerified`
	}
	query = query + `  WHERE id = @id`

	args := pgx.NamedArgs{
		"email":         auth.Email,
		"id":            id,
		"code":          auth.Code,
		"codeCreatedAt": auth.CodeCreatedAt,
		"isVerified":    auth.IsVerified,
		"updatedAt":     time.Now(),
	}

	_ = repo.connect.Db.QueryRow(ctx, query, args)

	return auth, nil
}

func NewRepository(db *storage.Postgres, logger *slog.Logger) domain.Repository {
	return &DbRepository{connect: db, logger: logger}
}
