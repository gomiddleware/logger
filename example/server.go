package main

import (
	"net/http"
	"os"
	"time"

	"github.com/go-kit/kit/log"
	gmLogger "github.com/gomiddleware/logger"
	"github.com/gomiddleware/realip"
	"github.com/gomiddleware/reqid"
)

func handler(w http.ResponseWriter, r *http.Request) {
	logger := gmLogger.LoggerFromRequest(r)
	logger.Log("evt", "handler.start")
	defer logger.Log("evt", "handler.end")

	w.Write([]byte(r.URL.Path + "\n"))
}

func utc() time.Time {
	return time.Now().UTC()
}

func main() {
	// create our own logger
	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stdout))
	logger = log.With(logger,
		"ts", log.TimestampFormat(utc, "20060102-150405.000000000"),
		"caller", log.DefaultCaller,
	)

	// create the logger middleware
	logMiddleware := gmLogger.New(logger)

	// Wrap each middleware in a chain (executed in reverse order, so ScrubRequestIdHeader then RandomId then RealIp
	// then logMiddleware then handler).
	middleware := logMiddleware(http.HandlerFunc(handler))
	middleware = realip.RealIp(middleware)
	middleware = reqid.RandomId(middleware)
	middleware = reqid.ScrubRequestIdHeader(middleware)

	// make the http.Hander and wrap it with the log middleware
	http.Handle("/", middleware)

	http.ListenAndServe(":8080", nil)
}
