package handler

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/peng225/orochi/internal/manager/api/server"
	"github.com/peng225/orochi/internal/manager/service"
)

type DatastoreHandler struct {
	dss *service.DatastoreService
}

func NewDatastoreHandler(dss *service.DatastoreService) *DatastoreHandler {
	return &DatastoreHandler{
		dss: dss,
	}
}

func (dsh *DatastoreHandler) CreateDatastore(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	data, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("Failed to read body.", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	var req server.CreateDatastoreRequest
	err = json.Unmarshal(data, &req)
	if err != nil {
		slog.Error("Failed to unmarshal body.", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	id, err := dsh.dss.CreateDatastore(r.Context(), *req.BaseURL)
	if err != nil {
		slog.Error("Failed to create datastore.", "err", err)
		switch {
		case errors.Is(err, service.ErrInvalidParameter):
			w.WriteHeader(http.StatusBadRequest)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	w.Header().Add("X-Datastore-ID", strconv.FormatInt(id, 10))
	w.WriteHeader(http.StatusCreated)
}

func (dsh *DatastoreHandler) GetDatastore(w http.ResponseWriter, r *http.Request, id server.DatastoreID) {
	ds, err := dsh.dss.GetDatastore(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrNotFound):
			w.WriteHeader(http.StatusNotFound)
		default:
			slog.Error("Failed to get datastore.", "err", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	data, err := json.Marshal(ds)
	if err != nil {
		slog.Error("Failed to marshal datastore.", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	_, err = w.Write(data)
	if err != nil {
		slog.Error("Failed to write data.", "err", err)
	}
}

func (dsh *DatastoreHandler) ListDatastores(
	w http.ResponseWriter, r *http.Request, params server.ListDatastoresParams,
) {
	var startFrom int64 = 0
	if params.StartFrom != nil {
		startFrom = *params.StartFrom
	}
	limit := 1000
	if params.Limit != nil {
		limit = *params.Limit
	}
	dsList, nextDatastoreID, err := dsh.dss.ListDatastores(r.Context(), startFrom, limit)
	if err != nil {
		slog.Error("Failed to list datastores.", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Add("X-Next-Datastore-ID", strconv.FormatInt(nextDatastoreID, 10))
	data, err := json.Marshal(dsList)
	if err != nil {
		slog.Error("Failed to marshal datastore list.", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = w.Write(data)
	if err != nil {
		slog.Error("Failed to write datastore list to body.", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
