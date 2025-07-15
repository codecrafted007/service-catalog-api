package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/codecrafted007/service-catalog-api/model"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

// mockStorage is a minimal in-memory mock of the Storage interface
type mockStorage struct {
	services []model.Service
	service  *model.Service
}

func (m *mockStorage) ListServices(context.Context, string, string, int, int) ([]model.Service, error) {
	return m.services, nil
}

func (m *mockStorage) CreateService(context.Context, *model.Service) (int64, error) {
	return 123, nil
}

func (m *mockStorage) DB() *sqlx.DB {
	return nil
}
func (m *mockStorage) IsValidAPIKey(key string) bool {
	return true
}

func (m *mockStorage) GetServiceById(xtx context.Context, id int) (*model.Service, error) {
	return m.service, nil
}
func (m *mockStorage) UpdateService(xtx context.Context, id int, s *model.Service) error {
	return nil
}
func (m *mockStorage) DeleteService(xtx context.Context, id int) error {
	return nil
}

func (m *mockStorage) CreateVersion(ctx context.Context, v *model.Version) (int64, error) {
	return 0, nil
}
func (m *mockStorage) GetVersionsByServiceID(ctx context.Context, serviceID int64) ([]*model.Version, error) {
	return nil, nil
}
func (m *mockStorage) GetVersionByID(ctx context.Context, versionID int64) (*model.Version, error) {
	return nil, nil
}
func (m *mockStorage) DeleteVersionByID(ctx context.Context, versionID int64) (bool, error) {
	return false, nil
}

func TestListServices(t *testing.T) {
	mock := &mockStorage{
		services: []model.Service{
			{ID: 1, Name: "Test Service", Description: "Mocked service"},
		},
	}

	logger := zap.NewNop().Sugar()
	h := NewServiceHandler(mock, logger)

	req := httptest.NewRequest(http.MethodGet, "/services?filter=Test&page=1&sort=name", nil)
	rec := httptest.NewRecorder()

	h.ListServices(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "Test Service")
}

func TestGetServiceByID(t *testing.T) {
	mock := &mockStorage{
		service: &model.Service{
			ID:          1,
			Name:        "Test Service",
			Description: "Mocked service",
		},
	}

	logger := zap.NewNop().Sugar()
	h := NewServiceHandler(mock, logger)

	r := mux.NewRouter()
	r.HandleFunc("/services/{id}", h.GetServiceByID).Methods("GET")
	req := httptest.NewRequest(http.MethodGet, "/services/1", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "Test Service")

}
