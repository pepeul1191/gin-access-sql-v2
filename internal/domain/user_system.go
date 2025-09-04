package domain

import "time"

type SystemUser struct {
	ID       uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	SystemID uint      `gorm:"not null" json:"system_id"`
	UserID   uint      `gorm:"not null" json:"user_id"`
	Created  time.Time `gorm:"not null" json:"created"`
}

func (SystemUser) TableName() string {
	return "systems_users"
}

type SystemUserItem struct {
	ID       int  `json:"id"`
	Selected bool `json:"selected"`
}
