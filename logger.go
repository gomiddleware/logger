package logger

import (
	"context"
	"net/http"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/gomiddleware/realip"
	"github.com/gomiddleware/reqid"
)

type key int

const loggerIdKey key = 82

// Logger middleware.
type Logger struct {
	h      http.Handler
	logger log.Logger
}

// SetLogger sets the logger to `log`. If you have used logger.New(), you can use this to set your
// logger. Alternatively, if you already have your log.Logger, then you can just call logger.NewLogger() directly.
func (l *Logger) SetLogger(logger log.Logger) {
	l.logger = logger
}

// wrapper to capture status.
type wrapper struct {
	http.ResponseWriter
	written int
	status  int
}

// capture status.
func (w *wrapper) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

// capture written bytes.
func (w *wrapper) Write(b []byte) (int, error) {
	n, err := w.ResponseWriter.Write(b)
	w.written += n
	return n, err
}

// NewLogger logger middleware with the given log.Logger.
func New(logger log.Logger) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return &Logger{
			logger: logger,
			h:      h,
		}
	}
}

// ServeHTTP implementation.
func (l *Logger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	res := &wrapper{w, 0, 200}

	// get the context since we'll use it a few times
	ctx := r.Context()

	// Get the RequestID (should have been set in previous middleware using gomiddleware/reqid)
	// and put it into the logger which goes in the context.
	rid := reqid.ReqIdFromContext(ctx)
	logger := log.With(l.logger, "rid", rid)

	// log the request.start
	logger.Log(
		"method", r.Method,
		"uri", r.RequestURI,
		"ip", realip.RealIpFromContext(ctx),
		"evt", "request.start",
	)

	// continue to the next middleware
	ctx = context.WithValue(r.Context(), loggerIdKey, logger)
	l.h.ServeHTTP(res, r.WithContext(ctx))

	// log the request.end
	logger.Log(
		"status", res.status,
		"size", res.written,
		"duration", time.Since(start),
		"evt", "request.end",
	)
}

// LoggerFromRequest can be used to obtain the Log from the request.
func LoggerFromRequest(r *http.Request) log.Logger {
	return r.Context().Value(loggerIdKey).(log.Logger)
}

// LoggerFromContext can be used to obtain the Log from the request.
func LoggerFromContext(ctx context.Context) log.Logger {
	return ctx.Value(loggerIdKey).(log.Logger)
}
