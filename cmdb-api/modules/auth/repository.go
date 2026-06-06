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

// GetOrCreateResourceType finds or creates a resource type by name
func (r *AuthRepository) GetOrCreateResourceType(name, description string) (*ResourceType, error) {
	var rt ResourceType
	err := database.DB.Where("name = ?", name).First(&rt).Error
	if err == nil {
		return &rt, nil
	}
	rt = ResourceType{Name: name, Description: description}
	if err := database.DB.Create(&rt).Error; err != nil {
		return nil, err
	}
	return &rt, nil
}

// GetOrCreatePermission finds or creates a permission
func (r *AuthRepository) GetOrCreatePermission(name string, resourceTypeID uint) (*Permission, error) {
	var p Permission
	err := database.DB.Where("name = ? AND resource_type_id = ?", name, resourceTypeID).First(&p).Error
	if err == nil {
		return &p, nil
	}
	p = Permission{Name: name, ResourceTypeID: resourceTypeID}
	if err := database.DB.Create(&p).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

// CreateResource creates a resource record
func (r *AuthRepository) CreateResource(res *Resource) error {
	return database.DB.Create(res).Error
}

// CheckRolePermission checks if a role has a specific permission on a resource
func (r *AuthRepository) CheckRolePermission(roleID uint, resourceID uint, permissionID uint) (bool, error) {
	var count int64
	err := database.DB.Model(&RolePermission{}).
		Where("role_id = ? AND resource_id = ? AND permission_id = ?", roleID, resourceID, permissionID).
		Count(&count).Error
	return count > 0, err
}

// GetUserPermittedResources returns resources of a specific type that the user has permission on
func (r *AuthRepository) GetUserPermittedResources(userID uint, resourceTypeName string, permissionName string) ([]Resource, error) {
	var resources []Resource
	err := database.DB.
		Joins("JOIN cmdb_auth.resource_types ON resource_types.id = resources.resource_type_id").
		Joins("JOIN cmdb_auth.role_permissions ON role_permissions.resource_id = resources.id").
		Joins("JOIN cmdb_auth.permissions ON permissions.id = role_permissions.permission_id").
		Joins("JOIN cmdb_auth.roles ON roles.id = role_permissions.role_id").
		Joins("JOIN cmdb_auth.user_roles ON user_roles.role_id = roles.id").
		Where("user_roles.user_id = ? AND resource_types.name = ? AND permissions.name = ?", userID, resourceTypeName, permissionName).
		Distinct().
		Find(&resources).Error
	return resources, err
}
