package accesslog

import (
	"net/http"
	"time"
)

type Level int

const (
	ResponseStatusSuccess     Level = 1 << iota // 200-299
	ResponseStatusRedirect                      // 300-399
	ResponseStatusClientError                   // 400-499
	ResponseStatusServerError                   // 500-599

	ResponseStatusNone Level = 0
	ResponseStatusAll        = ResponseStatusSuccess | ResponseStatusRedirect | ResponseStatusClientError | ResponseStatusServerError
)

func (l Level) InRange(statusCode int) bool {
	return (l&ResponseStatusSuccess != 0 && statusCode >= 200 && statusCode <= 299) ||
		(l&ResponseStatusRedirect != 0 && statusCode >= 300 && statusCode <= 399) ||
		(l&ResponseStatusClientError != 0 && statusCode >= 400 && statusCode <= 499) ||
		(l&ResponseStatusServerError != 0 && statusCode >= 500 && statusCode <= 599)
}

type AccessLogger interface {
	LogAccess(a Access)
}

type Access struct {
	// Request received from the client.
	Request *http.Request
	// The Time the request was received.
	Time time.Time
	// The Duration of handling the request.
	Duration time.Duration
	// StatusCode which was sent back to the client.
	StatusCode int
	// Size of the response body returned to the client.
	Size int
}

const panicNilHandler = "accesslog.Collect: http.Handler must not be nil"

func Collect(log AccessLogger, level Level, h *http.Handler) {
	if h == nil {
		panic(panicNilHandler)
	}
	if log == nil || level == ResponseStatusNone {
		return
	}

	*h = &collector{
		next:  *h,
		level: level,
		log:   log,
	}
}

type collector struct {
	next  http.Handler
	level Level
	log   AccessLogger
}

func (c *collector) ServeHTTP(wri http.ResponseWriter, req *http.Request) {
	start := time.Now()
	resp := responseWrapper{ResponseWriter: wri}
	c.next.ServeHTTP(&resp, req)

	sc := resp.StatusCode()
	if !c.level.InRange(sc) {
		return
	}

	c.log.LogAccess(Access{
		Request:    req,
		Time:       start,
		Duration:   time.Since(start),
		StatusCode: sc,
		Size:       resp.size,
	})
}

type responseWrapper struct {
	http.ResponseWriter
	statusCode int
	size       int
}

func (rw *responseWrapper) Write(b []byte) (int, error) {
	n, err := rw.ResponseWriter.Write(b)
	rw.size += n
	return n, err
}

func (rw *responseWrapper) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

func (rw *responseWrapper) StatusCode() int {
	if rw.statusCode != 0 {
		return rw.statusCode
	}
	return http.StatusOK
}
