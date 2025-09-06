package domain

import "time"

type SystemUserPermission struct {
	ID           uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	SystemID     uint      `gorm:"not null" json:"system_id"`
	UserID       uint      `gorm:"not null" json:"user_id"`
	PermissionID uint      `gorm:"not null" json:"permission_id"`
	Created      time.Time `gorm:"not null" json:"created"`
}

func (SystemUserPermission) TableName() string {
	return "systems_users_permissions"
}

type UserSystemPermission struct {
	SystemID       uint64 `json:"-"`
	SystemName     string `json:"-"`
	RoleID         uint64 `json:"-"`
	RoleName       string `json:"-"`
	PermissionID   uint64 `json:"-"`
	PermissionName string `json:"-"`
}
