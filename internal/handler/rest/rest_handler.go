package rest

import (
	"auth-service/internal/domain"
	"auth-service/internal/server"
	"encoding/json"
	"golang.org/x/exp/slog"
	"io"
	"net/http"
	"strconv"
	"time"
)

type handler struct {
	s      domain.Service
	logger *slog.Logger
}

type Handler interface {
	CreateAuth(w http.ResponseWriter, r *http.Request)
	FetchAuth(w http.ResponseWriter, r *http.Request)
	ResetPassword(w http.ResponseWriter, r *http.Request)
	DeleteAuth(w http.ResponseWriter, r *http.Request)
	UpdateAuth(w http.ResponseWriter, r *http.Request)
}

type createAuthResponse struct {
	ID            int       `json:"id"`
	Email         string    `json:"email"`
	Code          string    `json:"code"`
	CodeCreatedAt time.Time `json:"codeCreatedAt"`
	IsVerified    bool      `json:"is_verified"`
}

func (h *handler) CreateAuth(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var authRequest domain.AuthRepo
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			h.logger.Error("error while close:", err)
		}
	}(r.Body)

	err := json.NewDecoder(r.Body).Decode(&authRequest)

	if err != nil {
		h.logger.Error(err.Error())
		server.RespondWithError([]byte("error"), w, 400)
		return
	}
	createdAuth, err := h.s.CreateAuth(ctx, authRequest)

	if err != nil {
		h.logger.Error(err.Error())
		server.RespondWithError([]byte("server error"), w, 500)
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
		h.logger.Error(err.Error())
		server.RespondWithError([]byte("server error"), w, 500)
		return
	}

	server.RespondOK(jsonData, w)
}

func (h *handler) UpdateAuth(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			h.logger.Error("error while close:", err)
		}
	}(r.Body)

	data := struct {
		Email string `json:"email"`
		Id    string `json:"id"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&data)

	if err != nil {
		h.logger.Error("rest error while request decode: ", err.Error())
		server.RespondWithError([]byte("server error"), w, 500)
		return
	}

	updateAuth, err := h.s.UpdateAuth(ctx, domain.AuthUpdate{Email: data.Email, ID: data.Id})

	if err != nil {
		server.RespondWithError([]byte("No domain on this id"), w, 400)
		return
	}

	jsonData, err := json.Marshal(updateAuth)
	if err != nil {
		h.logger.Error(err.Error())
		server.RespondWithError([]byte("server error"), w, 500)
		return
	}
	server.RespondOK(jsonData, w)
}

func (h *handler) FetchAuth(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			h.logger.Error("error while close:", err)
		}
	}(r.Body)

	data := struct {
		Email string `json:"email"`
		Id    string `json:"id"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&data)

	if err != nil {
		h.logger.Error("rest error whilte request decode: ", err.Error())
		server.RespondWithError([]byte("server error"), w, 500)
		return
	}

	fetchAuth, err := h.s.FetchAuth(ctx, data.Id, data.Email)

	if err != nil {
		server.RespondWithError([]byte("No domain on this id"), w, 400)
		return
	}

	jsonData, err := json.Marshal(fetchAuth)
	if err != nil {
		h.logger.Error(err.Error())
		server.RespondWithError([]byte("server error"), w, 500)
		return
	}
	server.RespondOK(jsonData, w)
}

func (h *handler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var authRequest domain.AuthRepo

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			h.logger.Error("error while close:", err)
		}
	}(r.Body)

	err := json.NewDecoder(r.Body).Decode(&authRequest)

	if err != nil {
		h.logger.Error(err.Error())
		server.RespondWithError([]byte("error"), w, 500)
		return
	}
	isSend := h.s.ResetPassword(ctx, strconv.Itoa(authRequest.ID), authRequest.Email)

	server.RespondOK([]byte(strconv.FormatBool(isSend)), w)
}

func (h *handler) DeleteAuth(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			h.logger.Error("error while close:", err)
		}
	}(r.Body)

	id := r.URL.Query().Get("id")

	if id == "" {
		h.logger.Error("auth handler: empty id")
		server.RespondWithError([]byte("id empty"), w, 400)
		return
	}
	isOk := h.s.DeleteAuth(ctx, id)

	server.RespondOK([]byte(strconv.FormatBool(isOk)), w)
}

func NewHandler(s domain.Service, logger *slog.Logger) Handler {
	return &handler{s: s, logger: logger}
}
