package handler

import (
	"net/http"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (hh *HealthHandler) CheckHealthStatus(w http.ResponseWriter, r *http.Request) {
}
