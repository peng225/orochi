package handler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/peng225/orochi/internal/gateway/api/server"
	"github.com/peng225/orochi/internal/gateway/service"
)

type ObjectHandler struct {
	os *service.ObjectService
}

func NewObjectHandler(os *service.ObjectService) *ObjectHandler {
	return &ObjectHandler{
		os: os,
	}
}

func (oh *ObjectHandler) CreateObject(w http.ResponseWriter, r *http.Request, bucket server.Bucket, object server.Object) {
	defer r.Body.Close()
	err := oh.os.CreateObject(r.Context(), string(object), string(bucket), r.Body)
	if err != nil {
		slog.Error("CreateObject failed.", "err", err)
		switch {
		case errors.Is(err, service.ErrInvalidParameter):
			w.WriteHeader(http.StatusBadRequest)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (oh *ObjectHandler) GetObject(w http.ResponseWriter, r *http.Request, bucket server.Bucket, object server.Object) {
	data, err := oh.os.GetObject(r.Context(), string(object), string(bucket))
	if err != nil {
		slog.Error("GetObject failed.", "err", err)
		switch {
		case errors.Is(err, service.ErrInvalidParameter):
			w.WriteHeader(http.StatusBadRequest)
		case errors.Is(err, service.ErrNotFound):
			w.WriteHeader(http.StatusNotFound)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	_, err = w.Write(data)
	if err != nil {
		slog.Error("Failed to write data.", "err", err)
		return
	}
}

func (oh *ObjectHandler) DeleteObject(w http.ResponseWriter, r *http.Request, bucket server.Bucket, object server.Object) {
	err := oh.os.DeleteObject(r.Context(), string(object), string(bucket), r.Body)
	if err != nil {
		slog.Error("DeleteObject failed.", "err", err)
		switch {
		case errors.Is(err, service.ErrInvalidParameter):
			w.WriteHeader(http.StatusBadRequest)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (oh *ObjectHandler) ListObjects(
	w http.ResponseWriter, r *http.Request, bucket server.Bucket, params server.ListObjectsParams,
) {
	var startFrom int64 = 0
	if params.StartFrom != nil {
		startFrom = *params.StartFrom
	}
	limit := 1000
	if params.Limit != nil {
		limit = *params.Limit
	}
	objList, nextObjectID, err := oh.os.ListObjects(r.Context(), string(bucket), startFrom, limit)
	if err != nil {
		slog.Error("ListObjects failed.", "err", err)
		switch {
		case errors.Is(err, service.ErrInvalidParameter):
			w.WriteHeader(http.StatusBadRequest)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	w.Header().Add("X-Next-Object-ID", strconv.FormatInt(nextObjectID, 10))
	data, err := json.Marshal(objList)
	if err != nil {
		slog.Error("Failed to marshal object list.", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = w.Write(data)
	if err != nil {
		slog.Error("Failed to write object list to body.", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
