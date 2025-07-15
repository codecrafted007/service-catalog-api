package model

import "time"

type Version struct {
	ID        int64     `db:"id" json:"id"`
	ServiceID int64     `db:"service_id" json:"serviceId"`
	Version   string    `db:"version" json:"version"`
	Changelog string    `db:"changelog" json:"changelog,omitempty"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
}
