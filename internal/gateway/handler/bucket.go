package handler

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/peng225/orochi/internal/entity"
	"github.com/peng225/orochi/internal/gateway/api/private/client"
	"github.com/peng225/orochi/internal/gateway/api/private/server"
	"github.com/peng225/orochi/internal/gateway/service"
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
	var req client.CreateBucketRequest
	err = json.Unmarshal(data, &req)
	if err != nil {
		slog.Error("Failed to unmarshal data.", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = bh.bs.CreateBucket(*req.Name, *req.EcConfig)
	if err != nil {
		slog.Error("Failed to create bucket.", "err", err)
		if errors.Is(err, service.ErrInvalidParameter) {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
	w.WriteHeader(http.StatusCreated)
}

func (bh *BucketHandler) UpdateBucketStatusToDeleting(
	w http.ResponseWriter, r *http.Request, name server.BucketName,
) {
	err := bh.bs.UpdateBucketStatus(name, entity.BucketStatusDeleting)
	if err != nil {
		slog.Error("Failed to update bucket status.", "err", err)
		if errors.Is(err, service.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func (bh *BucketHandler) DeleteBucket(
	w http.ResponseWriter, r *http.Request, name server.BucketName,
) {
	bh.bs.DeleteBucket(name)
	w.WriteHeader(http.StatusNoContent)
}
