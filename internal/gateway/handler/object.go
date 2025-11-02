package handler

import (
	"errors"
	"log/slog"
	"net/http"

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

func (oh *ObjectHandler) CreateObject(w http.ResponseWriter, r *http.Request, bucket server.BucketParam, object server.ObjectParam) {
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

func (oh *ObjectHandler) GetObject(w http.ResponseWriter, r *http.Request, bucket server.BucketParam, object server.ObjectParam) {
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
