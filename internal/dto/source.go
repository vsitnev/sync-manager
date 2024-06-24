package dto

import (
	"time"
)

type Source struct {
	ID            *int    `json:"id" db:"id"`
	Name          string  `json:"name" db:"name"`
	Description   string  `json:"description" db:"description"`
	Code          string  `json:"code" db:"code"`
	ReceiveMethod string  `json:"receive_method" db:"receive_method"`
	Routes        []Route `json:"routes"`

	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" db:"updated_at"`
}

type Route struct {
	ID   *int   `json:"id,omitempty" db:"id"`
	Name string `json:"name" db:"name"`
	Url  string `json:"url" db:"url"`
}
