package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/codecrafted007/service-catalog-api/internal/storage"
	"github.com/codecrafted007/service-catalog-api/internal/utils"
	"github.com/codecrafted007/service-catalog-api/model"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type ServiceHandler struct {
	Store  storage.Storage
	Logger *zap.SugaredLogger
}

func NewServiceHandler(store storage.Storage, logger *zap.SugaredLogger) *ServiceHandler {
	return &ServiceHandler{
		Store:  store,
		Logger: logger,
	}
}

func (h *ServiceHandler) ListServices(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	filter := r.URL.Query().Get("filter")
	sort := r.URL.Query().Get("sort")

	pageStr := r.URL.Query().Get("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	limitStr := r.URL.Query().Get("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 20
	}

	services, err := h.Store.ListServices(ctx, filter, sort, page, limit)

	if err != nil {
		h.Logger.Error("error listing services: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, nil, "Internal Server Error")
		return
	}

	utils.WriteJSON(w, http.StatusOK, services, "")
}

func (h *ServiceHandler) GetServiceByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.Logger.Error("invalid service id: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, nil, "Invalid service ID")
		return
	}

	svc, err := h.Store.GetServiceById(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			h.Logger.Warnw("Service not found", "id", id)
			utils.WriteJSON(w, http.StatusNotFound, nil, "Service not found")
			return
		}

		h.Logger.Errorw("Failed to get service", "error", err)
		utils.WriteJSON(w, http.StatusInternalServerError, nil, "Internal server error")
		return
	}

	utils.WriteJSON(w, http.StatusOK, svc, "")
}

func (h *ServiceHandler) CreateService(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var input struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Version     string `json:"version"`
		Changelog   string `json:"changelog,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.Logger.Errorw("invalid input", "error", err)
		utils.WriteJSON(w, http.StatusBadRequest, nil, "Invalid input")
		return
	}

	service := model.Service{
		Name:        input.Name,
		Description: input.Description,
		CreatedAt:   time.Now(),
	}

	id, err := h.Store.CreateService(ctx, &service)
	if err != nil {
		h.Logger.Errorw("failed to create service", "error", err)
		utils.WriteJSON(w, http.StatusInternalServerError, nil, "Failed to create service")
		return
	}
	service.ID = int(id)

	version := model.Version{
		ServiceID: id,
		Version:   input.Version,
		Changelog: input.Changelog,
		CreatedAt: time.Now(),
	}

	_, err = h.Store.CreateVersion(ctx, &version)
	if err != nil {
		h.Logger.Errorw("failed to create initial version", "error", err)
		utils.WriteJSON(w, http.StatusInternalServerError, nil, "Service created but failed to add version")
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]int64{"id": id}, "")
}

func (h *ServiceHandler) UpdateService(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	serviceIDStr := mux.Vars(r)["id"]
	serviceID, err := strconv.Atoi(serviceIDStr)
	if err != nil {
		h.Logger.Errorw("invalid service ID", "error", err)
		utils.WriteJSON(w, http.StatusBadRequest, nil, "Invalid service ID")
		return
	}

	var input struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.Logger.Errorw("invalid input", "error", err)
		utils.WriteJSON(w, http.StatusBadRequest, nil, "Invalid input")
		return
	}

	updatedService := model.Service{
		Name:        input.Name,
		Description: input.Description,
	}

	err = h.Store.UpdateService(ctx, serviceID, &updatedService)
	if err != nil {
		h.Logger.Errorw("failed to update service", "error", err)
		utils.WriteJSON(w, http.StatusInternalServerError, nil, "Failed to update service")
		return
	}

	utils.WriteJSON(w, http.StatusOK, nil, "")
}

func (h *ServiceHandler) DeleteService(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, nil, "invalid service id")
		return
	}

	err = h.Store.DeleteService(ctx, id)
	if err != nil {
		h.Logger.Errorf("failed to delete service %d: %v", id, err)
		utils.WriteJSON(w, http.StatusInternalServerError, nil, "could not delete service")
		return
	}

	utils.WriteJSON(w, http.StatusOK, "service deleted successfully", "")
}
