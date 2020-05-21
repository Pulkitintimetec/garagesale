package handlers

import (
	"context"
	"net/http"

	"garagesale/007.errorhandling/internal/platform/database"
	"garagesale/007.errorhandling/internal/platform/web"
)

// Heartbeat HTTP handler type
type Heartbeat struct{}

// Health returns a HTTP status of OK
func (h *Heartbeat) Health(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var status struct {
		Status string
	}
	if err := database.Open(); err != nil {
		status.Status = "db not ready"
	}
	status.Status = "OK"
	web.Respond(ctx, w, status, http.StatusOK)
	return nil
}
