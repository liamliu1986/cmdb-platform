package auth

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Username  string         `gorm:"size:32;uniqueIndex;not null" json:"username"`
	Nickname  string         `gorm:"size:20" json:"nickname"`
	Email     string         `gorm:"size:100;uniqueIndex;not null" json:"email"`
	Password  string         `gorm:"size:80;not null" json:"-"`
	Status    int            `gorm:"default:1" json:"status"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (User) TableName() string { return "cmdb_auth.users" }

type Role struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"size:64;uniqueIndex;not null" json:"name"`
	Description string         `gorm:"size:255" json:"description"`
	IsAdmin     bool           `gorm:"default:false" json:"is_admin"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Role) TableName() string { return "cmdb_auth.roles" }

type RoleRelation struct {
	ID       uint `gorm:"primaryKey" json:"id"`
	ParentID uint `gorm:"not null;index" json:"parent_id"`
	ChildID  uint `gorm:"not null;index" json:"child_id"`
}

func (RoleRelation) TableName() string { return "cmdb_auth.role_relations" }

type ResourceType struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"size:64;uniqueIndex;not null" json:"name"`
	Description string    `gorm:"size:255" json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

func (ResourceType) TableName() string { return "cmdb_auth.resource_types" }

type Resource struct {
	ID             uint         `gorm:"primaryKey" json:"id"`
	Name           string       `gorm:"size:128;not null" json:"name"`
	ResourceTypeID uint         `gorm:"not null;index" json:"resource_type_id"`
	ResourceType   ResourceType `gorm:"foreignKey:ResourceTypeID" json:"resource_type,omitempty"`
	CreatedAt      time.Time    `json:"created_at"`
}

func (Resource) TableName() string { return "cmdb_auth.resources" }

type ResourceGroup struct {
	ID             uint         `gorm:"primaryKey" json:"id"`
	Name           string       `gorm:"size:64;not null" json:"name"`
	ResourceTypeID uint         `gorm:"not null" json:"resource_type_id"`
	ResourceType   ResourceType `gorm:"foreignKey:ResourceTypeID" json:"resource_type,omitempty"`
	CreatedAt      time.Time    `json:"created_at"`
}

func (ResourceGroup) TableName() string { return "cmdb_auth.resource_groups" }

type ResourceGroupItem struct {
	ID         uint `gorm:"primaryKey" json:"id"`
	GroupID    uint `gorm:"not null;index" json:"group_id"`
	ResourceID uint `gorm:"not null;index" json:"resource_id"`
}

func (ResourceGroupItem) TableName() string { return "cmdb_auth.resource_group_items" }

type Permission struct {
	ID             uint   `gorm:"primaryKey" json:"id"`
	Name           string `gorm:"size:32;not null" json:"name"`
	ResourceTypeID uint   `gorm:"not null" json:"resource_type_id"`
}

func (Permission) TableName() string { return "cmdb_auth.permissions" }

type RolePermission struct {
	ID           uint `gorm:"primaryKey" json:"id"`
	RoleID       uint `gorm:"not null;index" json:"role_id"`
	ResourceID   uint `gorm:"not null;index" json:"resource_id"`
	PermissionID uint `gorm:"not null;index" json:"permission_id"`
}

func (RolePermission) TableName() string { return "cmdb_auth.role_permissions" }

type UserRole struct {
	ID     uint `gorm:"primaryKey" json:"id"`
	UserID uint `gorm:"not null;index" json:"user_id"`
	RoleID uint `gorm:"not null;index" json:"role_id"`
}

func (UserRole) TableName() string { return "cmdb_auth.user_roles" }
