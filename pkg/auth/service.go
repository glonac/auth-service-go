package auth

type AuthService struct {
	repo *DbRepository
}

func (s *AuthService) CreateAuth(auth Auth) {
	s.repo.CreateAuth(auth)
}

func NewService(repo *DbRepository) *AuthService {
	return &AuthService{repo}
}
