package storage

import (
	"context"

	"github.com/codecrafted007/service-catalog-api/model"
	"github.com/jmoiron/sqlx"
)

type Storage interface {
	ListServices(ctx context.Context, filter string, sort string, page, limit int) ([]model.Service, error)
	GetServiceById(ctx context.Context, id int) (*model.Service, error)
	CreateService(ctx context.Context, s *model.Service) (int64, error)
	UpdateService(ctx context.Context, id int, s *model.Service) error
	DeleteService(ctx context.Context, id int) error

	IsValidAPIKey(key string) bool
	DB() *sqlx.DB

	CreateVersion(ctx context.Context, v *model.Version) (int64, error)
	GetVersionsByServiceID(ctx context.Context, serviceID int64) ([]*model.Version, error)
	GetVersionByID(ctx context.Context, versionID int64) (*model.Version, error)
	DeleteVersionByID(ctx context.Context, versionID int64) (bool, error)
}
