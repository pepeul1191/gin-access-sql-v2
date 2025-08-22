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
