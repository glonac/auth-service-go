package handler

// импорты не сортируются через goimports
// следует настроить IDE
import (
	"auth-service/internal/auth"
	"auth-service/internal/logger"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

var (
	// неканониченое название ошибки
	// советую ознакомиться с конвенциями вроде:
	// https://github.com/golang/go/wiki/CodeReviewComments
	// https://go.dev/doc/effective_go
	// https://github.com/uber-go/guide/blob/master/style.md
	ErrorHandle = errors.New("error while handle request")
)

type RestHandler struct {
	s auth.IAuthService
}

type createAuthResponse struct {
	ID            uint      `json:"id"`
	Email         string    `json:"email"`
	Code          string    `json:"code"`
	CodeCreatedAt time.Time `json:"codeCreatedAt"`
	IsVerified    bool      `json:"is_verified"`
}

func (h *RestHandler) createAuth(w http.ResponseWriter, r *http.Request) {
	var auth auth.Auth
	// The Server will close the request body. The ServeHTTP
	// Handler does not need to.
	// закрывать body нужно только в клиентах для респонсов
	defer r.Body.Close()

	err := json.NewDecoder(r.Body).Decode(&auth)

	if err != nil {
		// эта штука же просто вернет атрибут, она ничего не залогирует
		// получается логирования в хэндлерах вообще нет
		logger.Err(err)
		// все же обычно ошибку возвращают в json-объекте
		w.Write([]byte(ErrorHandle.Error()))
		return
	}
	createdAuth, err := h.s.CreateAuth(auth)

	if err != nil {
		logger.Err(err)
		w.Write([]byte(err.Error()))
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
		w.Write([]byte(ErrorHandle.Error()))
	}

	w.Write(jsonData)
}

func (h *RestHandler) FetchAuth(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	data := struct {
		Email string `json:"email"`
		Id    string `json:"id"`
	}{}
	err := json.NewDecoder(r.Body).Decode(&data)

	if err != nil {
		logger.Err(err)
	}
	auth := h.s.FetchAuth(data.Id, data.Email)

	if auth.ID == 0 {
		w.Write([]byte("No auth on this id"))
		return
	}

	jsonData, err := json.Marshal(auth)
	if err != nil {
		logger.Err(err)
	}

	w.Write(jsonData)
}

func (h *RestHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var auth auth.Auth
	defer r.Body.Close()

	err := json.NewDecoder(r.Body).Decode(&auth)

	if err != nil {
		log.Println(err)
	}

	// conversion from uint to string yields a string of one rune, not a string of digits (did you mean fmt.Sprint(x)?)
	isSend := h.s.ResetPassword(string(auth.ID), auth.Email)
	w.Write([]byte(strconv.FormatBool(isSend)))
}

// это скорее setup роутера, чем handle
func (h *RestHandler) HandleRequests(router *chi.Mux) {
	router.Route("/", func(r chi.Router) {
		r.Post("/sign-up", h.createAuth)
		r.Post("/fetch-user", h.FetchAuth)
		r.Post("/resetPass", h.ResetPassword)
		//TODO verify-email
	})
}

func NewHandler(s auth.IAuthService) *RestHandler {
	return &RestHandler{s: s}
}
