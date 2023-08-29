package server

import (
	"net/http"
)

func RespondWithError(msg []byte, w http.ResponseWriter, status int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_, _ = w.Write(msg)
	w.WriteHeader(status)
}
