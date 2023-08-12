package auth

import (
	"auth-service/internal/grpc/userGrpc"
	"auth-service/internal/logger"
	"auth-service/internal/user"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const separatorCode = "|"

// unused
const lifeCodeHours = 24

type authService struct {
	repo       AuthRepository
	userClient user.UserClient
}

// В go интерфейсы не именуют с заглавный буквы I — это не java/c#
type IAuthService interface {
	CreateAuth(auth Auth) (createdAuth Auth, err error)
	FetchAuth(id, email string) Auth
	ResetPassword(id, email string) bool
}

// unused
type encodeCode struct {
	email      string
	createTime string
}

var (
	// уже другой, но все еще неканоничный нейминг ошибки
	validationError = errors.New("validate error")
)

// большой красный флаг, что метод не принимает контекст и не передает его дальше
// таким образом если клиент прервет свой реквест то сервер будет ждать пока метод не отработает
func (s *authService) CreateAuth(auth Auth) (createdAuth Auth, err error) {
	if auth.Email == "" || auth.Password == "" {
		// чтобы понапрасну не создавать объекты обычно возвращают *Auth и соответствтенно nil
		return Auth{}, validationError
	}

	// непонятна семантика поля Code
	code := s.generateCode(auth.Email)
	auth.Code = code

	createAuth, err := s.repo.CreateAuth(auth)

	// хочется оборачивать ошибки в подобных кейсах
	if err != nil {
		return Auth{}, err
	}

	//TODO add logic to send email notification with code

	//TODO add logic to check exist user or create new
	s.userClient.GetUserByClientAccountId(context.Background(), &userGrpc.GetUserByClientAccountId_Request{
		ClientAccountId: strconv.Itoa(int(createAuth.ID))})

	return createAuth, nil
}

// метод не возвращает ошибку
func (s *authService) FetchAuth(id, email string) Auth {
	auth, err := s.repo.FetchAuth(id, email)
	if err != nil {
		logger.Err(err)
	}
	return auth
}

// тут вместо bool-а нужно возвращать ошибку
func (s *authService) ResetPassword(id, email string) bool {
	auth, err := s.repo.FetchAuth(id, email)
	if err != nil {
		logger.Err(err)
		return false
	}
	// log
	fmt.Println(auth)
	//TODO add logic to send in notification service code to reset pass
	return false
}

// TODO test this
func (s *authService) generateCode(email string) string {
	currentTime := time.Now()
	//TODO change on solid hash algorithm
	return base64.StdEncoding.EncodeToString([]byte(email + separatorCode + currentTime.String()))
}

// TODO test this
// unused
func (s *authService) encodeCode(code string) (codeStruct encodeCode, isValid bool) {
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

// очень редко вы хотите на самом деле возвращать интерфейс
// обычно интерфейс принимают, а возврвщают конкретную реализацию
func NewService(
	repo AuthRepository,
	grpcUserClient user.UserClient,
) IAuthService {
	return &authService{repo: repo, userClient: grpcUserClient}
}
