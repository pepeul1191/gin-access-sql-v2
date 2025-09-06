// internal/domain/user.go
package domain

import "time"

type Role struct {
	ID          uint         `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string       `gorm:"size:20;not null" json:"name"`
	Created     time.Time    `gorm:"not null" json:"created"`
	Updated     time.Time    `gorm:"not null" json:"updated"`
	SystemID    uint         `gorm:"not null" json:"system_id"`                   // Hace referencia a System.ID
	System      System       `gorm:"foreignKey:SystemID" json:"system,omitempty"` // Relación con System
	Permissions []Permission `json:"permissions"`
}
