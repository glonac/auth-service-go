package rest_test

import (
	"auth-service/internal/handler/rest"
	"auth-service/internal/repositories"
	"auth-service/mocks"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestRestHandler_FetchAuth(t *testing.T) {
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
			mockError:    errors.New("No auth on this id"),
		},
		{
			name:         "Search by id",
			id:           "1",
			mockResId:    "1",
			mockResEmail: "test@mail.ru",
		},
		{
			name:         "Empty all struct",
			email:        "",
			id:           "",
			mockResId:    "0",
			mockResEmail: "",
			mockError:    errors.New("No auth on this id"),
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mockAuthRepo := mocks.NewAuthService(t)
			idUuid, _ := strconv.ParseUint(tc.mockResId, 10, 32)
			mockAuthRepo.On("FetchAuth", tc.id, tc.email).
				Return(repositories.Auth{Email: tc.mockResEmail, ID: uint(idUuid)}, tc.mockError).
				Once()
			handler := rest.NewHandler(mockAuthRepo)

			input := fmt.Sprintf(`{"id": "%s", "email": "%s"}`, tc.id, tc.email)
			req, err := http.NewRequest(http.MethodPost, "/fetch-user", bytes.NewReader([]byte(input)))
			require.NoError(t, err)
			rr := httptest.NewRecorder()

			handler.FetchAuth(rr, req)

			require.Equal(t, rr.Code, http.StatusOK)

			body := rr.Body.String()
			var expectedRes struct {
				ID    uint
				Email string
			}
			if tc.mockError != nil {
				require.EqualError(t, tc.mockError, body)
			} else {
				json.Unmarshal([]byte(body), &expectedRes)
				require.Equal(t, struct {
					ID    uint
					Email string
				}{
					ID:    uint(idUuid),
					Email: tc.mockResEmail,
				},
					expectedRes)
			}
		})
	}
}

func TestRestHandler_HandleRequests(t *testing.T) {

}
