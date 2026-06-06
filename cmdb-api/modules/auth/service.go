package auth

import (
	"errors"

	"golang.org/x/crypto/bcrypt"

	"cmdb-api/config"
	"cmdb-api/database"
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
	return s.CheckResourcePermission(userID, resourceName, 0, permissionName)
}

// CheckResourcePermission checks if user has permission on a specific resource
func (s *AuthService) CheckResourcePermission(userID uint, resourceTypeName string, resourceID uint, permissionName string) bool {
	roles, _ := s.repo.GetUserRoles(userID)
	for _, role := range roles {
		if role.IsAdmin {
			return true
		}
		if resourceID == 0 {
			continue
		}
		hasPerm, _ := s.repo.CheckRolePermission(role.ID, resourceID, s.getPermissionID(resourceTypeName, permissionName))
		if hasPerm {
			return true
		}
	}
	return false
}

// getPermissionID gets permission ID by resource type and permission name
func (s *AuthService) getPermissionID(resourceTypeName, permissionName string) uint {
	var rt ResourceType
	database.DB.Where("name = ?", resourceTypeName).First(&rt)
	var p Permission
	database.DB.Where("name = ? AND resource_type_id = ?", permissionName, rt.ID).First(&p)
	return p.ID
}

// InitIAMResources initializes default IAM resource types and permissions for Subnet
func (s *AuthService) InitIAMResources() error {
	rt, err := s.repo.GetOrCreateResourceType("Subnet", "IPAM Subnet")
	if err != nil {
		return err
	}
	perms := []string{"subnet:view", "subnet:allocate", "subnet:release", "subnet:admin"}
	for _, p := range perms {
		if _, err := s.repo.GetOrCreatePermission(p, rt.ID); err != nil {
			return err
		}
	}
	return nil
}

// GetUserPermittedSubnets returns subnet IDs that the user has allocate permission on
func (s *AuthService) GetUserPermittedSubnets(userID uint) ([]uint, error) {
	roles, _ := s.repo.GetUserRoles(userID)
	for _, role := range roles {
		if role.IsAdmin {
			return nil, nil // admin can access all, return nil to signal "all"
		}
	}
	resources, err := s.repo.GetUserPermittedResources(userID, "Subnet", "subnet:allocate")
	if err != nil {
		return nil, err
	}
	var ids []uint
	for _, r := range resources {
		ids = append(ids, r.ID)
	}
	return ids, nil
}
