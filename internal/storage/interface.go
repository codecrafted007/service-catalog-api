package storage

import (
	"github.com/codecrafted007/service-catalog-api/model"
	"github.com/jmoiron/sqlx"
)

type Storage interface {
	ListServices(filter string, sort string, page, limit int) ([]model.Service, error)
	GetServiceById(id int) (*model.Service, error)
	IsValidAPIKey(key string) bool
	DB() *sqlx.DB

	CreateService(s *model.Service) (int64, error)
	UpdateService(id int, s *model.Service) error
	DeleteService(id int) error
}
