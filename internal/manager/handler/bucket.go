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
	}
	var req server.CreateBucketRequest
	err = json.Unmarshal(data, &req)
	if err != nil {
		slog.Error("Failed to unmarshal body.", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	id, err := bh.bs.CreateBucket(r.Context(), *req.Name)
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
