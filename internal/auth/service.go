package auth

import (
	"auth-service/internal/grpc/client"
	"auth-service/internal/grpc/userGrpc"
	"auth-service/internal/logger"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const separatorCode = "|"
const lifeCodeHours = 24

type AuthService struct {
	repo           AuthRepository
	grpcUserClient *client.UserClientGrpc
}

type IAuthService interface {
	CreateAuth(auth Auth) (createdAuth Auth, err error)
	FetchAuth(id, email string) Auth
	ResetPassword(id, email string) bool
}

type encodeCode struct {
	email      string
	createTime string
}

var (
	validationError = errors.New("validate error")
)

func (s *AuthService) CreateAuth(auth Auth) (createdAuth Auth, err error) {
	if auth.Email == "" || auth.Password == "" {
		return Auth{}, validationError
	}

	code := s.generateCode(auth.Email)
	auth.Code = code

	createAuth, err := s.repo.CreateAuth(auth)

	if err != nil {
		return Auth{}, err
	}

	//TODO add logic to send email notification with code

	//TODO add logic to check exist user or create new
	s.grpcUserClient.Client.GetUserById(context.Background(), &userGrpc.GetUserById_Request{
		UserId: strconv.Itoa(int(createAuth.ID))})

	return createAuth, nil
}

func (s *AuthService) FetchAuth(id, email string) Auth {
	auth, err := s.repo.FetchAuth(id, email)
	if err != nil {
		logger.Err(err)
	}
	return auth
}

func (s *AuthService) ResetPassword(id, email string) bool {
	auth, err := s.repo.FetchAuth(id, email)
	if err != nil {
		logger.Err(err)
		return false
	}
	fmt.Println(auth)
	//TODO add logic to send in notification service code to reset pass
	return false
}

// TODO test this
func (s *AuthService) generateCode(email string) string {
	currentTime := time.Now()
	//TODO change on solid hash algorithm
	return base64.StdEncoding.EncodeToString([]byte(email + separatorCode + currentTime.String()))
}

// TODO test this
func (s *AuthService) encodeCode(code string) (codeStruct encodeCode, isValid bool) {
	codeString, err := base64.StdEncoding.DecodeString(code)
	if err != nil {
		logger.Err(err)
		return codeStruct, false
	}
	splitCode := strings.Split(string(codeString), separatorCode)

	if len(splitCode) != 2 {
		return codeStruct, false
	}
	currentTime := time.Now()
	createdTime, err := time.Parse("", splitCode[1])
	if err != nil {
		return codeStruct, false
	}

	lifeTimeEnd := createdTime.Add(lifeCodeHours * time.Hour)
	if currentTime.After(lifeTimeEnd) {
		return codeStruct, false
	}

	return encodeCode{
		email:      splitCode[0],
		createTime: splitCode[1],
	}, true
}

func NewService(
	repo AuthRepository,
	grpcUserClient *client.UserClientGrpc,
) IAuthService {
	return &AuthService{repo: repo, grpcUserClient: grpcUserClient}
}
