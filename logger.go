package logger

import (
	"net/http"
	"os"
	"time"

	"github.com/chilts/sid"
	"github.com/gomiddleware/logit"
)

// Logger middleware.
type Logger struct {
	h   http.Handler
	log *logit.Logger
}

// SetLogger sets the logger to `log`. If you have used logger.New(), you can use this to set your
// logger. Alternatively, if you already have your log.Logger, then you can just call logger.NewLogger() directly.
func (l *Logger) SetLogger(log *logit.Logger) {
	l.log = log
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

// New logger middleware.
func New() func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return &Logger{
			log: logit.New(os.Stdout, "req"),
			h:   h,
		}
	}
}

// NewLogger logger middleware with the given logger.
func NewLogger(log *logit.Logger) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return &Logger{
			log: log.Clone("req"),
			h:   h,
		}
	}
}

// ServeHTTP implementation.
func (l *Logger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	res := &wrapper{w, 0, 200}

	// ToDo: we could use `r.Header.Get("X-Request-ID")` but we should have a config value to determine whether to use
	// it or not, since it could come from outside and be untrusted.

	// set up some fields for this request
	l.log.WithField("method", r.Method)
	l.log.WithField("uri", r.RequestURI)
	l.log.WithField("id", sid.Id())
	l.log.Log("request-start")

	// continue to the next middleware
	l.h.ServeHTTP(res, r)

	// output the final log line
	l.log.WithField("status", res.status)
	l.log.WithField("size", res.written)
	l.log.WithField("duration", time.Since(start))
	l.log.Log("request-end")
}
