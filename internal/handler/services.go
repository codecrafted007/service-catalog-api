package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/codecrafted007/service-catalog-api/internal/storage"
	"github.com/codecrafted007/service-catalog-api/internal/utils"
	"github.com/codecrafted007/service-catalog-api/model"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type Handler struct {
	Store  storage.Storage
	Logger *zap.SugaredLogger
}

func New(store storage.Storage, logger *zap.SugaredLogger) *Handler {
	return &Handler{
		Store:  store,
		Logger: logger,
	}
}

func (h *Handler) ListServices(w http.ResponseWriter, r *http.Request) {
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

	services, err := h.Store.ListServices(filter, sort, page, limit)

	if err != nil {
		h.Logger.Error("error listing services: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, nil, "Internal Server Error")
		return
	}

	utils.WriteJSON(w, http.StatusOK, services, "")
}

func (h *Handler) GetServiceByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.Logger.Error("invalid service id: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, nil, "Invalid service ID")
		return
	}

	svc, err := h.Store.GetServiceById(id)
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

func (h *Handler) CreateService(w http.ResponseWriter, r *http.Request) {
	var svc model.Service
	if err := json.NewDecoder(r.Body).Decode(&svc); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, nil, "invalid request")
		return
	}

	if svc.Name == "" {
		utils.WriteJSON(w, http.StatusBadRequest, nil, "service name is mandatory")
		return
	}

	id, err := h.Store.CreateService(&svc)
	if err != nil {
		h.Logger.Errorf("failed to create service: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, nil, "could not create service")
		return
	}

	utils.WriteJSON(w, http.StatusOK, id, "")
}

func (h *Handler) UpdateService(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, nil, "invalid service id")
		return
	}

	var svc model.Service
	if err := json.NewDecoder(r.Body).Decode(&svc); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, nil, "invalid request body")
		return
	}

	if svc.Name == "" {
		utils.WriteJSON(w, http.StatusBadRequest, nil, "service name is required")
		return
	}
	err = h.Store.UpdateService(id, &svc)
	if err != nil {
		h.Logger.Errorf("failed to update service %d: %v", id, err)
		utils.WriteJSON(w, http.StatusInternalServerError, nil, "could not update service")
		return
	}

	utils.WriteJSON(w, http.StatusOK, "service updated", "")
}

func (h *Handler) DeleteService(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, nil, "invalid service id")
		return
	}

	err = h.Store.DeleteService(id)
	if err != nil {
		h.Logger.Errorf("failed to delete service %d: %v", id, err)
		utils.WriteJSON(w, http.StatusInternalServerError, nil, "could not delete service")
		return
	}

	utils.WriteJSON(w, http.StatusOK, "service deleted successfully", "")
}
