package domain

import (
	"auth-service/internal/grpc/client"
	"context"
	"encoding/base64"
	"fmt"
	"golang.org/x/exp/slog"
	"strings"
	"time"
)

const separatorCode = "|"
const lifeCodeHours = 24

type authService struct {
	log            *slog.Logger
	repo           Repository
	grpcUserClient *client.UserClientGrpc
	//queueClient    queue.QueueService
}

type AuthRepo struct {
	ID            int
	Password      string
	Email         string
	Code          string
	CodeCreatedAt time.Time
	IsVerified    bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type AuthFetch struct {
	ID            int
	Email         string
	Code          string
	CodeCreatedAt time.Time
	IsVerified    bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type AuthUpdate struct {
	ID            string
	Email         string
	Code          string
	CodeCreatedAt time.Time
	IsVerified    bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type encodeCode struct {
	email      string
	createTime string
}

func (s *authService) CreateAuth(ctx context.Context, auth AuthRepo) (createdAuth AuthFetch, err error) {
	if auth.Email == "" || auth.Password == "" {
		return AuthFetch{}, ValidationError
	}

	code := s.generateCode(auth.Email)
	auth.Code = code

	createAuth, err := s.repo.CreateAuth(ctx, auth)
	fmt.Println(createAuth)
	if err != nil {
		s.log.Error("authService: " + err.Error())
		return AuthFetch{}, err
	}

	//TODO add logic to send email notification with code

	//TODO add logic to check exist user or create new
	//_, err = s.grpcUserClient.Client.GetUserById(context.Background(), &userGrpc.GetUserById_Request{
	//	UserId: strconv.Itoa(int(createAuth.ID))})
	//
	//if err != nil {
	//	s.log.Error("authService: %s", err)
	//	return AuthRepo{}, err
	//}
	//err = s.queueClient.SendMsg("test")
	//if err != nil {
	//	s.log.Error("authService: %s", err)
	//	return AuthRepo{}, err
	//}
	return createAuth, nil
}

func (s *authService) UpdateAuth(ctx context.Context, auth AuthUpdate) (AuthUpdate, error) {
	auth, err := s.repo.UpdateAuth(ctx, auth.ID, auth)
	if err != nil {
		s.log.Error("auth service:", err.Error())
		return auth, err
	}
	return auth, nil
}

func (s *authService) FetchAuth(ctx context.Context, id, email string) (AuthFetch, error) {
	auth, err := s.repo.FetchAuth(ctx, id, email)
	if err != nil {
		s.log.Error("auth service:", err.Error())
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

func (s *authService) DeleteAuth(ctx context.Context, id string) bool {
	isOk := s.repo.DeleteAuth(ctx, id)
	return isOk
}

func NewService(
	repo Repository,
	grpcUserClient *client.UserClientGrpc,
	logger *slog.Logger,
	// queue queue.QueueService,
) Service {
	return &authService{repo: repo, grpcUserClient: grpcUserClient, log: logger} //queueClient: queue

}
