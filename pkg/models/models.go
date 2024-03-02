package models

import "time"

type Project struct {
	ID        string    `json:"id,omitempty"`
	Name      string    `json:"name,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

type Good struct {
	ID          int    	  `json:"id"`
	ProjectID   int   	  `json:"project"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Priority    int       `json:"priority"`
	Removed     bool      `json:"removed"`
	CreatedAt   time.Time `json:"created_at"`
}

type CreateGoodRequest struct {
	Name string `json:"name" validate:"required"`
}

type ReprioritizeGoodRequest struct {
	Priority int `json:"newPriority" validate:"required"`
}

type ReprioritizeGoodResponse struct {
	ID        string `json:"id"`
	Priority    int    `json:"priority"`
}

type Meta struct {
	Limit   int `json:"limit"`
	Ofset   int `json:"offset"`
	Total   int `json:"total"`
	Removed int `json:"removed"`
}

type GoodsResponse struct {
	Meta  Meta	`json:"meta"`
	Goods []Good `json:"goods"`
}
