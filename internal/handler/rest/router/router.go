package router

import (
	"auth-service/internal/handler/rest"
	"github.com/go-chi/chi/v5"
)

type router struct {
	h rest.Handler
}

type Router interface {
	HandleRequests(r *chi.Mux)
}

func (h *router) HandleRequests(r *chi.Mux) {
	r.Post("/sign-up", h.h.CreateAuth)
	r.Post("/fetch-user", h.h.FetchAuth)
	r.Post("/update-user", h.h.UpdateAuth)
	r.Post("/reset-pass", h.h.ResetPassword)
	r.Delete("/delete", h.h.DeleteAuth)
	//TODO verify-email
}

func NewRouter(h rest.Handler) Router {
	return &router{h: h}
}
