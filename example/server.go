package main

import (
	"net/http"
	"os"

	"github.com/gomiddleware/logger"
	"github.com/gomiddleware/logit"
)

func handler(logger *logit.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Log("inside handler")
		w.Write([]byte(r.URL.Path))
	}
}

func main() {
	// create the logger middleware
	lggr := logit.New(os.Stdout, "main")
	log := logger.NewLogger(lggr)

	// make the http.Hander and wrap it with the log middleware
	handle := http.HandlerFunc(handler(lggr))
	http.Handle("/", log(handle))

	http.ListenAndServe(":8080", nil)
}
