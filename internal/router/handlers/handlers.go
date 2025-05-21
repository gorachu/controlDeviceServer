package handlers

import (
	"controlDeviceServer/internal/storage/sqlite"
	"log/slog"
)

type Handler struct {
	storage *sqlite.Storage
	log     *slog.Logger
}

func NewHandlers(storage *sqlite.Storage, log *slog.Logger) *Handler {
	return &Handler{storage: storage, log: log}
}
