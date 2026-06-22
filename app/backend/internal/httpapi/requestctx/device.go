package requestctx

import (
	"context"
	"net/http"

	"github.com/ajn2004/daleego-hearth/backend/internal/db"
)

type contextKey string

const currentDeviceKey contextKey = "currentDevice"

func WithCurrentDevice(ctx context.Context, device db.Device) context.Context {
	return context.WithValue(ctx, currentDeviceKey, device)
}

func CurrentDevice(r *http.Request) (db.Device, bool) {
	device, ok := r.Context().Value(currentDeviceKey).(db.Device)
	return device, ok
}
