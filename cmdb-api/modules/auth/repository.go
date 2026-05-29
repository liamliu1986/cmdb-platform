package auth

import (
	"cmdb-api/database"
)

type AuthRepository struct{}

func NewAuthRepository() *AuthRepository {
	return &AuthRepository{}
}

// User
func (r *AuthRepository) CreateUser(user *User) error {
	return database.DB.Create(user).Error
}

func (r *AuthRepository) GetUserByUsername(username string) (*User, error) {
	var user User
	err := database.DB.Where("username = ?", username).First(&user).Error
	return &user, err
}

func (r *AuthRepository) GetUserByID(id uint) (*User, error) {
	var user User
	err := database.DB.First(&user, id).Error
	return &user, err
}

func (r *AuthRepository) ListUsers(page, pageSize int) ([]User, int64, error) {
	var users []User
	var total int64
	db := database.DB.Model(&User{})
	db.Count(&total)
	err := db.Offset((page - 1) * pageSize).Limit(pageSize).Find(&users).Error
	return users, total, err
}

// Role
func (r *AuthRepository) CreateRole(role *Role) error {
	return database.DB.Create(role).Error
}

func (r *AuthRepository) GetRoleByID(id uint) (*Role, error) {
	var role Role
	err := database.DB.First(&role, id).Error
	return &role, err
}

func (r *AuthRepository) ListRoles() ([]Role, error) {
	var roles []Role
	err := database.DB.Find(&roles).Error
	return roles, err
}

// Permission
func (r *AuthRepository) GrantPermission(rp *RolePermission) error {
	return database.DB.Create(rp).Error
}

func (r *AuthRepository) GetRolePermissions(roleID uint) ([]RolePermission, error) {
	var rps []RolePermission
	err := database.DB.Where("role_id = ?", roleID).Find(&rps).Error
	return rps, err
}

// UserRole
func (r *AuthRepository) AssignRoleToUser(userID, roleID uint) error {
	return database.DB.Create(&UserRole{UserID: userID, RoleID: roleID}).Error
}

func (r *AuthRepository) GetUserRoles(userID uint) ([]Role, error) {
	var roles []Role
	err := database.DB.
		Joins("JOIN cmdb_auth.user_roles ON user_roles.role_id = roles.id").
		Where("user_roles.user_id = ?", userID).
		Find(&roles).Error
	return roles, err
}

// CheckRolePermission checks if a role has a specific permission on a resource
func (r *AuthRepository) CheckRolePermission(roleID uint, resourceName string, permissionName string) (bool, error) {
	var count int64
	err := database.DB.Table("cmdb_auth.role_permissions").
		Joins("JOIN cmdb_auth.resources ON resources.id = role_permissions.resource_id").
		Joins("JOIN cmdb_auth.permissions ON permissions.id = role_permissions.permission_id").
		Where("role_permissions.role_id = ?", roleID).
		Where("resources.name = ?", resourceName).
		Where("permissions.name = ?", permissionName).
		Count(&count).Error
	return count > 0, err
}
