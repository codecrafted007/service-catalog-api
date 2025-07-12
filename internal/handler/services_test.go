package handler

import (
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

func (m *mockStorage) ListServices(string, string, int, int) ([]model.Service, error) {
	return m.services, nil
}

func (m *mockStorage) CreateService(*model.Service) (int64, error) {
	return 123, nil
}

func (m *mockStorage) DB() *sqlx.DB {
	return nil
}
func (m *mockStorage) IsValidAPIKey(key string) bool {
	return true
}

func (m *mockStorage) GetServiceById(id int) (*model.Service, error) {
	return m.service, nil
}
func (m *mockStorage) UpdateService(id int, s *model.Service) error {
	return nil
}
func (m *mockStorage) DeleteService(id int) error {
	return nil
}

func TestListServices(t *testing.T) {
	mock := &mockStorage{
		services: []model.Service{
			{ID: 1, Name: "Test Service", Description: "Mocked service"},
		},
	}

	logger := zap.NewNop().Sugar()
	h := New(mock, logger)

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
	h := New(mock, logger)

	r := mux.NewRouter()
	r.HandleFunc("/services/{id}", h.GetServiceByID).Methods("GET")
	req := httptest.NewRequest(http.MethodGet, "/services/1", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "Test Service")

}
