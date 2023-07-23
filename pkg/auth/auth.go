package auth

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"log"
	"strconv"
	"time"
)

var (
	ErrorNoAuth     = errors.New("no such user")
	ErrorWhileFetch = errors.New("error while fetch")
)

type DbRepository struct {
	connect *gorm.DB
}

type Auth struct {
	gorm.Model
	ID            uint      `json:"id" gorm:"primary_key"`
	Password      string    `json:"password"`
	Email         string    `json:"email"`
	Code          string    `json:"code"`
	CodeCreatedAt time.Time `json:"codeCreatedAt"`
	IsVerified    bool      `json:"is_verified"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdateAt      time.Time `json:"updateAt"`
}

func NewAuth(password, email, code string, codeCreateAt time.Time, isVerified bool) Auth {
	return Auth{
		Password:      password,
		Email:         email,
		Code:          code,
		CodeCreatedAt: codeCreateAt,
		IsVerified:    isVerified,
	}
}

func (repo *DbRepository) CreateAuth(auth Auth) Auth {
	repo.connect.Create(&auth)
	return auth
}

func (repo *DbRepository) DeleteAuth(id string) {
	repo.connect.Delete(&Auth{}, id)
}

func (repo *DbRepository) FetchAuth(id, email string) (auth Auth, err error) {
	if id != "" {
		idUuid, err := strconv.ParseUint(id, 10, 32)
		if err != nil {
			return auth, fmt.Errorf("Error while transform id %s", err)
		}
		res := repo.connect.Find(auth, "id = ?", uint(idUuid))
		if res.Error != nil {
			log.Println(res.Error)
			return Auth{}, ErrorNoAuth
		}
		return auth, nil
	}
	if email != "" {
		auth := Auth{Email: email}
		res := repo.connect.Find(&auth, "email = ?", email)
		if res.Error != nil {
			log.Println(res.Error)
			return Auth{}, ErrorNoAuth
		}
		return auth, nil
	}
	return Auth{}, ErrorWhileFetch
}

func NewRepository(db *gorm.DB) *DbRepository {
	return &DbRepository{connect: db}
}
