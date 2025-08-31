// internal/domain/user.go
package domain

import "time"

type Permission struct {
	ID      uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name    string    `gorm:"size:20;not null" json:"name"`
	Created time.Time `gorm:"not null" json:"created"`
	Updated time.Time `gorm:"not null" json:"updated"`
	RoleID  uint      `gorm:"not null" json:"role"`                    // Hace referencia a System.ID
	Role    Role      `gorm:"foreignKey:RoleID" json:"role,omitempty"` // Relación con System
}
