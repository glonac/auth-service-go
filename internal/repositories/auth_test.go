package repositories_test

import (
	authDomain "auth-service/internal/domain"
	"auth-service/mocks"
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func TestDbRepository_CreateAuth(t *testing.T) {
	cases := []struct {
		name      string
		password  string
		email     string
		mockError error
	}{
		{
			name:     "Success",
			password: "check",
			email:    "test@mail.ru",
		},
		{
			name:     "Fail empty pass",
			password: "",
			email:    "test@mail.ru",
			//mockError: authDomain.ErrorWhileCreate,
		},
		{
			name:     "Fail empty email",
			password: "check",
			email:    "",
			//mockError: authDomain.ErrorWhileCreate,
		},
		{
			name:     "Fail invalid email",
			password: "check",
			email:    "test",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mockAuthRepo := mocks.NewAuthRepository(t)
			ctx := context.Background()
			if tc.mockError != nil {
				testStruct := authDomain.AuthRepo{Password: tc.password, Email: tc.email}
				mockAuthRepo.On("CreateAuth", testStruct).
					Return(authDomain.AuthRepo{Email: tc.email, Password: tc.password}, tc.mockError).
					Once()
				res, err := mockAuthRepo.CreateAuth(ctx, testStruct)
				fmt.Println(res)
				fmt.Println(tc)
				if err != nil {
					assert.EqualError(t, tc.mockError, err.Error())
				}
				assert.Equal(t, testStruct, res)
			}
		})
	}
}

func TestDbRepository_FetchAuth(t *testing.T) {
	cases := []struct {
		name         string
		email        string
		id           string
		mockError    error
		mockResId    string
		mockResEmail string
	}{
		{
			name:         "Success",
			email:        "test@mail.ru",
			id:           "1",
			mockResId:    "1",
			mockResEmail: "test@mail.ru",
		},
		{
			name:         "Fail broke error",
			email:        "notExist@email.ru",
			id:           "",
			mockResId:    "0",
			mockResEmail: "notExist@email.ru",
			mockError:    authDomain.ErrorNoAuth,
		},
		{
			name: "Search by id",
			id:   "1",
		},
		{
			name:         "Empty all struct",
			email:        "",
			id:           "",
			mockResId:    "0",
			mockResEmail: "",
			mockError:    authDomain.ErrorWhileFetch,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mockAuthRepo := mocks.NewAuthRepository(t)
			idUuid, _ := strconv.ParseUint(tc.mockResId, 10, 32)
			ctx := context.Background()
			if tc.mockError != nil {
				mockAuthRepo.On("FetchAuth", tc.id, tc.email).
					Return(authDomain.AuthRepo{Email: tc.mockResEmail, ID: uint(idUuid)}, tc.mockError).
					Once()
				res, err := mockAuthRepo.FetchAuth(ctx, tc.id, tc.email)
				if err != nil {
					assert.EqualError(t, tc.mockError, err.Error())
				}
				assert.Equal(t,
					authDomain.AuthRepo{
						Email: tc.mockResEmail,
						ID:    uint(idUuid)}, res)
			}
		})
	}
}
