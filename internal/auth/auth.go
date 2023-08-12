package auth

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"gorm.io/gorm"
)

var (
	ErrorNoAuth      = errors.New("no such user")
	ErrorWhileFetch  = errors.New("error while fetch")
	ErrorWhileCreate = errors.New("error while create")
)

type DbRepository struct {
	connect *gorm.DB
}

type AuthRepository interface {
	CreateAuth(auth Auth) (createAuth Auth, err error)
	DeleteAuth(id string)
	FetchAuth(id, email string) (auth Auth, err error)
	UpdateAuth(id string, auth Auth) (Auth, error)
}

type Auth struct {
	gorm.Model
	ID            uint      `json:"id" gorm:"primary_key"`
	Password      string    `json:"password"`
	Email         string    `json:"email"`
	Code          string    `json:"code"`
	CodeCreatedAt time.Time `json:"codeCreatedAt"` // разный кейс в json-тегах для разных полей
	IsVerified    bool      `json:"is_verified"`   // следует выбрать один
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

func (repo *DbRepository) CreateAuth(auth Auth) (createAuth Auth, err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
			err = ErrorWhileCreate
		}
	}()
	//TODO constraint on email model
	//TODO add logic to hash pass
	repo.connect.Create(&auth)
	return auth, nil
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
		fmt.Println(res)
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

func (repo *DbRepository) UpdateAuth(id string, auth Auth) (Auth, error) {
	res := repo.connect.Model(&auth).Where("id = ?", id).Updates(auth)
	if res.RowsAffected == 0 {
		return auth, errors.New("no update for this query")
	}
	return auth, nil
}

func NewRepository(db *gorm.DB) AuthRepository {
	return &DbRepository{connect: db}
}
