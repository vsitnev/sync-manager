package model

import (
	"database/sql"
	"github.com/vsitnev/sync-manager/internal/dto"
	"time"
)

type Source struct {
	ID            int     `json:"id" db:"id"`
	Name          string  `json:"name" db:"name"`
	Description   string  `json:"description" db:"description"`
	Code          string  `json:"code" db:"code"`
	ReceiveMethod string  `json:"receive_method" db:"receive_method"`
	Routes        []Route `json:"routes"`

	CreatedAt time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt sql.NullTime `json:"updated_at" db:"updated_at"`
}

func (m *Source) ToDto() dto.Source {
	var updated *time.Time
	if m.UpdatedAt.Valid {
		updated = &m.UpdatedAt.Time
	} else {
		updated = nil
	}

	var routes []dto.Route
	if m.Routes != nil && len(m.Routes) != 0 {
		for _, item := range m.Routes {
			routes = append(routes, dto.Route{
				Name: item.Name,
				Url:  item.Url,
			})
		}
	}

	return dto.Source{
		ID:            &m.ID,
		Name:          m.Name,
		Description:   m.Description,
		Code:          m.Code,
		ReceiveMethod: m.ReceiveMethod,
		Routes:        routes,
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     updated,
	}
}

type Route struct {
	ID       int    `json:"id" db:"id"`
	Name     string `json:"name" db:"name"`
	Url      string `json:"url" db:"url"`
	SourceID string `json:"source_id" db:"source_fk"`
}
