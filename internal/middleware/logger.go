package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		start := time.Now()
		wrapped := &statusWriter{ResponseWriter: writer, status: http.StatusOK}
		next.ServeHTTP(wrapped, req)

		duration := time.Since(start)
		attrs := []any{
			slog.String("method", req.Method),
			slog.String("path", req.URL.Path),
			slog.Int64("duration_ms", duration.Milliseconds()),
			slog.Duration("duration", duration),
		}

		if wrapped.status >= 500 {
			slog.Error("request failed", attrs...)
		} else if wrapped.status >= 400 {
			slog.Warn("client error", attrs...)
		} else {
			slog.Info("request completed", attrs...)
		}
	})
}
