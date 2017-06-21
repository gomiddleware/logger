// Package logger provides a logger which can be used as Go Middleware. It logs a line on an incoming request and
// another line when that request has finished, along with the final status code and the time elapsed.
//
// It also adds a RequestID to the logger which is pulled from the request context, which should be set with the
// gomiddleware/reqid middleware before this middleware runs.
//
// Incoming fields logged are : method, uri, ip.
//
// Outgoing fields logged are : status, size (bytes), duration.
//
package logger
