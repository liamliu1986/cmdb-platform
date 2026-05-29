package auth

import (
	"errors"

	"golang.org/x/crypto/bcrypt"

	"cmdb-api/config"
	"cmdb-api/pkg/jwtutil"
)

type AuthService struct {
	repo *AuthRepository
	cfg  *config.Config
}

func NewAuthService(cfg *config.Config) *AuthService {
	return &AuthService{repo: NewAuthRepository(), cfg: cfg}
}

func (s *AuthService) Register(req *RegisterRequest) (*User, error) {
	existing, _ := s.repo.GetUserByUsername(req.Username)
	if existing != nil && existing.ID > 0 {
		return nil, errors.New("username already exists")
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user := &User{
		Username: req.Username,
		Password: string(hashed),
		Nickname: req.Nickname,
		Email:    req.Email,
	}
	if err := s.repo.CreateUser(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *AuthService) Login(req *LoginRequest) (*LoginResponse, error) {
	user, err := s.repo.GetUserByUsername(req.Username)
	if err != nil {
		return nil, errors.New("invalid username or password")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid username or password")
	}
	token, err := jwtutil.GenerateToken(user.ID, user.Username, s.cfg.JWTSecret, s.cfg.JWTExpireHours)
	if err != nil {
		return nil, err
	}
	return &LoginResponse{Token: token, UserID: user.ID, Username: user.Username}, nil
}

func (s *AuthService) CheckPermission(userID uint, resourceName string, permissionName string) bool {
	roles, _ := s.repo.GetUserRoles(userID)
	for _, role := range roles {
		if role.IsAdmin {
			return true
		}
	}
	return false
}
