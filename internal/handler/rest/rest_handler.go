package rest

import (
	"auth-service/internal/auth"
	"auth-service/internal/logger"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"golang.org/x/exp/slog"
	"net/http"
	"strconv"
	"time"
)

var (
	ErrorHandle = errors.New("error while handle request")
)

type Handler struct {
	s      auth.AuthService
	logger *slog.Logger
}

type createAuthResponse struct {
	ID            uint      `json:"id"`
	Email         string    `json:"email"`
	Code          string    `json:"code"`
	CodeCreatedAt time.Time `json:"codeCreatedAt"`
	IsVerified    bool      `json:"is_verified"`
}

func (h *Handler) createAuth(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var authRequest auth.AuthRepoStruct
	defer r.Body.Close()

	err := json.NewDecoder(r.Body).Decode(&authRequest)

	if err != nil {
		h.logger.Error(err.Error())
		w.WriteHeader(400)
		_, _ = w.Write([]byte(ErrorHandle.Error()))
		return
	}
	createdAuth, err := h.s.CreateAuth(ctx, authRequest)

	if err != nil {
		h.logger.Error(err.Error())
		w.WriteHeader(500)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	response := createAuthResponse{
		ID:            createdAuth.ID,
		Email:         createdAuth.Email,
		Code:          createdAuth.Code,
		CodeCreatedAt: createdAuth.CodeCreatedAt,
		IsVerified:    createdAuth.IsVerified,
	}

	jsonData, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(500)
		_, _ = w.Write([]byte(ErrorHandle.Error()))
		return
	}
	w.WriteHeader(200)
	_, _ = w.Write(jsonData)
}

func (h *Handler) FetchAuth(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	defer r.Body.Close()
	data := struct {
		Email string `json:"email"`
		Id    string `json:"id"`
	}{}
	err := json.NewDecoder(r.Body).Decode(&data)

	if err != nil {
		h.logger.Error(err.Error())
	}
	fetchAuth, err := h.s.FetchAuth(ctx, data.Id, data.Email)
	if err != nil {
		w.WriteHeader(400)
		_, _ = w.Write([]byte("No auth on this id"))
		return
	}

	jsonData, err := json.Marshal(fetchAuth)
	if err != nil {
		logger.Err(err)
		w.WriteHeader(500)
		_, _ = w.Write([]byte("error"))
	}
	w.WriteHeader(200)
	_, _ = w.Write(jsonData)
}

func (h *Handler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var authRequest auth.AuthRepoStruct
	defer r.Body.Close()

	err := json.NewDecoder(r.Body).Decode(&authRequest)

	if err != nil {
		h.logger.Error(err.Error())
	}
	isSend := h.s.ResetPassword(ctx, string(authRequest.ID), authRequest.Email)
	w.WriteHeader(200)
	_, _ = w.Write([]byte(strconv.FormatBool(isSend)))
}

func (h *Handler) HandleRequests(r *chi.Mux) {
	r.Post("/sign-up", h.createAuth)
	r.Post("/fetch-user", h.FetchAuth)
	r.Post("/resetPass", h.ResetPassword)
	//TODO verify-email
}

func NewHandler(s auth.AuthService, logger *slog.Logger) *Handler {
	return &Handler{s: s, logger: logger}
}
