package model

import "time"

type Version struct {
	ID        int       `json:"id"`
	Version   string    `json:"version"`
	CreatedAt time.Time `json:"createdAt"`
}

type Service struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	Versions    []Version `json:"versions,omitempty"`
}
