package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/codecrafted007/service-catalog-api/internal/storage"
	"github.com/codecrafted007/service-catalog-api/internal/utils"
	"github.com/codecrafted007/service-catalog-api/model"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type VersionHandler struct {
	Store  storage.Storage
	Logger *zap.SugaredLogger
}

func NewVersionHandler(store storage.Storage, logger *zap.SugaredLogger) *VersionHandler {
	return &VersionHandler{
		Store:  store,
		Logger: logger,
	}
}

func (h *VersionHandler) CreateVersion(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	serviceIDStr := mux.Vars(r)["id"]
	serviceID, err := strconv.ParseInt(serviceIDStr, 10, 64)
	if err != nil {
		h.Logger.Warnw("Invalid service ID", "service_id", serviceIDStr, "error", err)
		utils.WriteJSON(w, http.StatusBadRequest, nil, "Invalid service ID")
		return
	}

	var input struct {
		Version   string `json:"version"`
		Changelog string `json:"changelog,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.Logger.Warnw("Failed to decode CreateVersion payload", "error", err)
		utils.WriteJSON(w, http.StatusBadRequest, nil, "Invalid input")
		return
	}

	newVersion := model.Version{
		ServiceID: serviceID,
		Version:   input.Version,
		Changelog: input.Changelog,
		CreatedAt: time.Now(),
	}

	insertedID, err := h.Store.CreateVersion(ctx, &newVersion)
	if err != nil {
		h.Logger.Errorw("DB error creating version", "service_id", serviceID, "error", err)
		utils.WriteJSON(w, http.StatusInternalServerError, nil, "Failed to create version")
		return
	}
	newVersion.ID = insertedID
	h.Logger.Infow("Version created successfully", "version_id", insertedID, "service_id", serviceID)
	utils.WriteJSON(w, http.StatusOK, newVersion, "")
}

// GET /services/{id}/versions
func (h *VersionHandler) ListVersions(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	serviceIDStr := mux.Vars(r)["id"]
	serviceID, err := strconv.ParseInt(serviceIDStr, 10, 64)
	if err != nil {
		h.Logger.Warnw("Invalid service ID", "service_id", serviceIDStr, "error", err)
		utils.WriteJSON(w, http.StatusBadRequest, nil, "Invalid service ID")
		return
	}

	versions, err := h.Store.GetVersionsByServiceID(ctx, serviceID)
	if err != nil {
		h.Logger.Errorw("Error while fetching versions", "service_id", serviceID, "error", err)
		utils.WriteJSON(w, http.StatusInternalServerError, nil, "Failed to fetch versions")
		return
	}
	h.Logger.Infow("Versions fetched succesfully", "service_id", serviceID, "count", len(versions))
	utils.WriteJSON(w, http.StatusOK, versions, "")
}

// GET /versions/{id}
func (h *VersionHandler) GetVersion(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	versionIDStr := mux.Vars(r)["id"]
	versionID, err := strconv.ParseInt(versionIDStr, 10, 64)
	if err != nil {
		h.Logger.Warnw("Invalid version ID", "version_id", versionIDStr, "error", err)
		utils.WriteJSON(w, http.StatusBadRequest, nil, "Invalid version ID")
		return
	}

	version, err := h.Store.GetVersionByID(ctx, versionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.Logger.Warnw("Version not found", "version_id", versionID)
			utils.WriteJSON(w, http.StatusNotFound, nil, "Version not found")
		} else {
			h.Logger.Errorw("Error while fetching version", "version_id", versionID, "error", err)
			utils.WriteJSON(w, http.StatusInternalServerError, nil, "Failed to fetch version")
		}
		return
	}
	h.Logger.Infow("Version fetched successfully", "version_id", versionID)
	utils.WriteJSON(w, http.StatusOK, version, "")
}

// DELETE /versions/{id}
func (h *VersionHandler) DeleteVersion(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	versionIDStr := mux.Vars(r)["id"]
	versionID, err := strconv.ParseInt(versionIDStr, 10, 64)
	if err != nil {
		h.Logger.Warnw("Invalid version ID for delete", "version_id", versionIDStr, "error", err)
		utils.WriteJSON(w, http.StatusBadRequest, nil, "Invalid version ID")
		return
	}

	deleted, err := h.Store.DeleteVersionByID(ctx, versionID)
	if err != nil {
		h.Logger.Errorw("Error while deleting version", "version_id", versionID, "error", err)
		utils.WriteJSON(w, http.StatusInternalServerError, nil, "Failed to delete version")
		return
	}

	if !deleted {
		h.Logger.Warnw("Version not found for deletion", "version_id", versionID)
		utils.WriteJSON(w, http.StatusNotFound, nil, "Version not found")
		return
	}
	h.Logger.Infow("Version deleted succesfully", "version_id", versionID)
	utils.WriteJSON(w, http.StatusNoContent, nil, "")
}
