package handler

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/peng225/orochi/internal/datastore/api/server"
	"github.com/peng225/orochi/internal/datastore/service"
)

type ObjectHandler struct {
	os *service.ObjectService
}

func NewObjectHandler(os *service.ObjectService) *ObjectHandler {
	return &ObjectHandler{
		os: os,
	}
}

func (oh *ObjectHandler) CreateObject(w http.ResponseWriter, r *http.Request, object server.Object) {
	defer r.Body.Close()
	err := oh.os.CreateObject(string(object), r.Body)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidParameter):
			w.WriteHeader(http.StatusBadRequest)
		case errors.Is(err, service.ErrTooLargeObject):
			w.WriteHeader(http.StatusRequestEntityTooLarge)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (oh *ObjectHandler) GetObject(w http.ResponseWriter, r *http.Request, object server.Object) {
	data, err := oh.os.GetObject(string(object))
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidParameter):
			w.WriteHeader(http.StatusBadRequest)
		case errors.Is(err, service.ErrObjectNotFound):
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

func (oh *ObjectHandler) DeleteObject(w http.ResponseWriter, r *http.Request, object server.Object) {
	err := oh.os.DeleteObject(string(object))
	if err != nil {
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
