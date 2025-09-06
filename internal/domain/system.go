package domain

import "time"

type System struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string    `gorm:"size:40;not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	Repository  string    `gorm:"size:100" json:"repository"`
	Created     time.Time `gorm:"not null" json:"created"`
	Updated     time.Time `gorm:"not null" json:"updated"`
	Roles       []*Role   `json:"roles"`
}

func (System) TableName() string {
	return "systems"
}

type SystemListResponse struct {
	Systems     []System `json:"systems"`
	Total       int64    `json:"total"`
	Page        int      `json:"page"`
	PerPage     int      `json:"per_page"`
	TotalPages  int      `json:"total_pages"`
	NameQuery   string   `json:"name_query,omitempty"`
	DescQuery   string   `json:"desc_query,omitempty"`
	StartRecord int      `json:"start_record"`
	EndRecord   int      `json:"end_record"`
}
