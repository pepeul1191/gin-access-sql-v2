package domain

import "time"

type System struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string    `gorm:"size:40;not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	Repository  string    `gorm:"size:100" json:"repository"`
	Created     time.Time `gorm:"not null" json:"created"`
	Updated     time.Time `gorm:"not null" json:"updated"`
}

func (System) TableName() string {
	return "systems"
}
