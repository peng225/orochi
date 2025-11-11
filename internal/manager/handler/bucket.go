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

type BucketHandler struct {
	bs *service.BucketService
}

func NewBucketHandler(bs *service.BucketService) *BucketHandler {
	return &BucketHandler{
		bs: bs,
	}
}

func (bh *BucketHandler) CreateBucket(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	data, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("Failed to read body.", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var req server.CreateBucketRequest
	err = json.Unmarshal(data, &req)
	if err != nil {
		slog.Error("Failed to unmarshal body.", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if req.Name == nil {
		slog.Error("Name is not set.")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if req.EcConfig == nil {
		slog.Error("EC config is not set.")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id, err := bh.bs.CreateBucket(r.Context(), *req.Name, *req.EcConfig)
	if err != nil {
		slog.Error("Failed to create bucket.", "err", err)
		switch {
		case errors.Is(err, service.ErrInvalidParameter):
			w.WriteHeader(http.StatusBadRequest)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	w.Header().Add("X-Bucket-ID", strconv.FormatInt(id, 10))
	w.WriteHeader(http.StatusCreated)
}

func (bh *BucketHandler) GetBucket(w http.ResponseWriter, r *http.Request, id int64) {
	bucket, err := bh.bs.GetBucket(r.Context(), id)
	if err != nil {
		slog.Error("Failed to get bucket.", "err", err)
		switch {
		case errors.Is(err, service.ErrNotFound):
			w.WriteHeader(http.StatusNotFound)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	data, err := json.Marshal(&bucket)
	if err != nil {
		slog.Error("Failed to marshal bucket.", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	_, err = w.Write(data)
	if err != nil {
		slog.Error("Failed to write data.", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (bh *BucketHandler) DeleteBucket(w http.ResponseWriter, r *http.Request, id int64) {
	err := bh.bs.DeleteBucket(r.Context(), id)
	if err != nil {
		slog.Error("Failed to delete bucket.", "err", err)
		switch {
		case errors.Is(err, service.ErrInvalidParameter):
			w.WriteHeader(http.StatusBadRequest)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusAccepted)
}
