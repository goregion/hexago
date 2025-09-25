// Package log provides multi-handler implementation for slog.
// This file contains the multiHandler type that enables logging to multiple
// destinations simultaneously, such as console, files, and remote services.
package log

import (
	"context"
	"log/slog"
)

// multiHandler implements slog.Handler interface to support multiple handlers simultaneously.
// This allows logging to multiple destinations (e.g., console and file) with a single logger.
// All handlers are called for each log record, providing comprehensive logging coverage.
type multiHandler struct {
	handlers []slog.Handler
}

// Enabled checks if logging is enabled for the given level in any of the handlers.
// Returns true if at least one handler would process a record at the specified level.
// This is used by slog to optimize performance by avoiding unnecessary work.
func (h *multiHandler) Enabled(ctx context.Context, level slog.Level) bool {
	for _, handler := range h.handlers {
		if handler.Enabled(ctx, level) {
			return true
		}
	}
	return false
}

// Handle processes a log record by forwarding it to all configured handlers.
// If any handler returns an error, the first error encountered is returned.
// This ensures that logging continues even if one handler fails, while still
// reporting errors for debugging purposes.
func (h *multiHandler) Handle(ctx context.Context, record slog.Record) error {
	var firstErr error
	for _, handler := range h.handlers {
		if err := handler.Handle(ctx, record); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}

// WithAttrs returns a new multiHandler where each underlying handler has the specified
// attributes added. This maintains the multi-handler structure while propagating
// the attribute addition to all wrapped handlers.
//
// This method is called when logger.With() is used to add structured attributes.
func (h *multiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	handlers := make([]slog.Handler, len(h.handlers))
	for i, handler := range h.handlers {
		handlers[i] = handler.WithAttrs(attrs)
	}
	return &multiHandler{
		handlers: handlers,
	}
}

// WithGroup returns a new multiHandler where each underlying handler has been
// configured with the specified group name. This maintains the multi-handler
// structure while propagating the group configuration to all wrapped handlers.
//
// This method is called when logger.WithGroup() is used to create hierarchical log structure.
func (h *multiHandler) WithGroup(name string) slog.Handler {
	handlers := make([]slog.Handler, len(h.handlers))
	for i, handler := range h.handlers {
		handlers[i] = handler.WithGroup(name)
	}
	return &multiHandler{
		handlers: handlers,
	}
}
