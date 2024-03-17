package util

import (
	"context"
	"log/slog"

	"github.com/taxio/errors"
)

func LogError(ctx context.Context, err error) {
	if err == nil {
		return
	}

	var cErr *errors.Error
	var attrs []any
	if errors.As(err, &cErr) {
		attrs = []any{
			slog.Any("attributes", cErr.Attributes()),
			slog.Any("stackTrace", ErrorStackFramePaths(err)),
		}
	}
	slog.ErrorContext(
		ctx, err.Error(),
		attrs...,
	)
}
