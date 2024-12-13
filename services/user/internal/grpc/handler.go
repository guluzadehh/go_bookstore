package grpchandler

import "log/slog"

type Handler struct {
	Log *slog.Logger
}

func NewHandler(log *slog.Logger) *Handler {
	return &Handler{
		Log: log,
	}
}
