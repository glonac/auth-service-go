package handler

import (
	"auth-service/pkg/auth"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"time"
)

type RestHandler struct {
	s *auth.AuthService
}

type createAuthRequest struct {
	Password      string    `json:"password"`
	Email         string    `json:"email"`
	Code          string    `json:"code"`
	CodeCreatedAt time.Time `json:"codeCreatedAt"`
	IsVerified    bool      `json:"is_verified"`
}

func (h *RestHandler) createAuth(w http.ResponseWriter, r *http.Request) {
	var auth auth.Auth
	defer r.Body.Close()

	err := json.NewDecoder(r.Body).Decode(&auth)

	if err != nil {
		log.Println(err)
	}
	h.s.CreateAuth(auth)
}

func (h *RestHandler) HandleRequests(router *chi.Mux) {
	router.Route("/", func(r chi.Router) {
		r.Post("/sign-up", h.createAuth)
	})
}

func NewHandler(s *auth.AuthService) *RestHandler {
	return &RestHandler{s: s}
}
