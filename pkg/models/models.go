package models

import "time"

type Project struct {
	ID        string    `json:"id,omitempty"`
	Name      string    `json:"name,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

type Goods struct {
	ID          string    `json:"id,omitempty"`
	Project     Project   `json:"project,omitempty"`
	Name        string    `json:"name,omitempty"`
	Description string    `json:"description,omitempty"`
	Priority    int       `json:"priority,omitempty"`
	Removed     bool      `json:"removed,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
}
