package forms

import (
	"time"
)

type SystemCreateInput struct {
	Name        string    `form:"name" binding:"required"`
	Description string    `form:"description"`
	Repository  string    `form:"repository"`
	Created     time.Time `form:"created" time_format:"2006-01-02"`
	Updated     time.Time `form:"updated" time_format:"2006-01-02"`
}

type SystemEditInput struct {
	Name        string `form:"name" binding:"required"`
	Description string `form:"description"`
	Repository  string `form:"repository"`
}
