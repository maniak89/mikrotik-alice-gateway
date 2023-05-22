package storage

import (
	"context"
	"errors"

	"mikrotik-alice-gateway/internal/models/storage"
)

var ErrInvalidState = errors.New("invalid state")

type Storage interface {
	Routers(ctx context.Context) ([]*storage.Router, error)
	Log(ctx context.Context, routerID string, level storage.LogLevel, msg string)
	UpdateHost(ctx context.Context, host *storage.Host) error
}
