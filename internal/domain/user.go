// internal/domain/user.go
package domain

import "time"

type User struct {
	ID            uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Username      string    `gorm:"size:20;not null" json:"username"`
	Password      string    `gorm:"size:100;not null" json:"password"`
	ActivationKey string    `gorm:"size:30" json:"activation_key,omitempty"`
	ResetKey      string    `gorm:"size:30" json:"reset_key,omitempty"`
	Email         string    `gorm:"size:50;unique;not null" json:"email"`
	Activated     bool      `gorm:"not null;default:false" json:"activated"`
	Created       time.Time `gorm:"not null" json:"created"`
	Updated       time.Time `gorm:"not null" json:"updated"`
}

type UserSummary struct {
	ID                uint   `gorm:"column:id" json:"id"`
	Username          string `gorm:"column:username" json:"username"`
	Email             string `gorm:"column:email" json:"email"`
	Activated         bool   `gorm:"column:activated" json:"activated"`
	AssociationStatus int    `gorm:"column:association_status" json:"association_status"`
}

type SystemUserRolesPermissions struct {
	UserID         uint   `json:"user_id"`
	SystemID       uint   `json:"system_id"`
	PermissionID   uint   `json:"permission_id"`
	PermissionName string `json:"permission_name"`
	RoleID         uint   `json:"role_id"`
	RoleName       string `json:"role_name"`
	IsAssigned     bool   `json:"is_assigned"`
}

// Permission represents a permission within a role.
type UserPermission struct {
	ID         uint   `json:"id"`
	Name       string `json:"name"`
	IsAssigned bool   `json:"is_assigned"`
}

// RoleWithPermissions represents a role and its associated permissions.
type RoleWithPermissions struct {
	ID          uint             `json:"id"`
	Name        string           `json:"name"`
	RoleID      string           `json:"role_id"`
	Permissions []UserPermission `json:"permissions"`
}
