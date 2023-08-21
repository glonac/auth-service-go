package auth

import (
	"auth-service/internal/grpc/client"
	"auth-service/internal/queue"
	"auth-service/pkg/grpc/userGrpc"
	"context"
	"encoding/base64"
	"errors"
	"golang.org/x/exp/slog"
	"gorm.io/gorm"
	"strconv"
	"strings"
	"time"
)

const separatorCode = "|"
const lifeCodeHours = 24

type authService struct {
	log            *slog.Logger
	repo           AuthRepository
	grpcUserClient *client.UserClientGrpc
	queueClient    queue.QueueService
}

type AuthService interface {
	CreateAuth(ctx context.Context, auth AuthRepoStruct) (createdAuth AuthRepoStruct, err error)
	FetchAuth(ctx context.Context, id, email string) (AuthRepoStruct, error)
	ResetPassword(ctx context.Context, id, email string) bool
}

type AuthRepository interface {
	CreateAuth(ctx context.Context, auth AuthRepoStruct) (createAuth AuthRepoStruct, err error)
	DeleteAuth(ctx context.Context, id string)
	FetchAuth(ctx context.Context, id, email string) (auth AuthRepoStruct, err error)
	UpdateAuth(ctx context.Context, id string, auth AuthRepoStruct) (AuthRepoStruct, error)
}

type AuthRepoStruct struct {
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

func (AuthRepoStruct) TableName() string {
	return "auth"
}

type encodeCode struct {
	email      string
	createTime string
}

var (
	ValidationError  = errors.New("validate error")
	ErrorNoAuth      = errors.New("no such user")
	ErrorWhileFetch  = errors.New("error while fetch")
	ErrorWhileCreate = errors.New("error while create")
)

func (s *authService) CreateAuth(ctx context.Context, auth AuthRepoStruct) (createdAuth AuthRepoStruct, err error) {
	if auth.Email == "" || auth.Password == "" {
		return AuthRepoStruct{}, ValidationError
	}

	code := s.generateCode(auth.Email)
	auth.Code = code

	createAuth, err := s.repo.CreateAuth(ctx, auth)

	if err != nil {
		s.log.Error("authService: " + err.Error())
		return AuthRepoStruct{}, err
	}

	//TODO add logic to send email notification with code

	//TODO add logic to check exist user or create new
	s.grpcUserClient.Client.GetUserById(context.Background(), &userGrpc.GetUserById_Request{
		UserId: strconv.Itoa(int(createAuth.ID))})
	s.queueClient.SendMsg("test")
	return createAuth, nil
}

func (s *authService) FetchAuth(ctx context.Context, id, email string) (AuthRepoStruct, error) {
	auth, err := s.repo.FetchAuth(ctx, id, email)
	if err != nil {
		s.log.Error(err.Error())
		return auth, err
	}
	return auth, nil
}

func (s *authService) ResetPassword(ctx context.Context, id, email string) bool {
	_, err := s.repo.FetchAuth(ctx, id, email)
	if err != nil {
		s.log.Error(err.Error())
		return false
	}
	//TODO add logic to send in notification service code to reset pass
	return false
}

func (s *authService) generateCode(email string) string {
	currentTime := time.Now()
	//TODO change on solid hash algorithm
	return base64.StdEncoding.EncodeToString([]byte(email + separatorCode + currentTime.String()))
}

func (s *authService) encodeCode(code string) (codeStruct encodeCode, isValid bool) {
	codeString, err := base64.StdEncoding.DecodeString(code)
	if err != nil {
		s.log.Error(err.Error())
		return codeStruct, false
	}
	splitCode := strings.Split(string(codeString), separatorCode)

	if len(splitCode) != 2 {
		return codeStruct, false
	}
	currentTime := time.Now()
	createdTime, err := time.Parse("", splitCode[1])
	if err != nil {
		s.log.Error(err.Error())
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
	logger *slog.Logger,
	queue queue.QueueService,
) AuthService {
	return &authService{repo: repo, grpcUserClient: grpcUserClient, log: logger, queueClient: queue}
}
