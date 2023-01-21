package accesslog

import (
	"github.com/go-pogo/serv/collect"
	"net"
	"net/http"
)

type Entry struct {
	Request *http.Request
	Metrics collect.Metrics
}

func RemoteAddr(r *http.Request) string {
	addr, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		addr = r.RemoteAddr
	}
	return addr
}

func (e *Entry) RemoteAddr() string { return RemoteAddr(e.Request) }

// https://www.rfc-editor.org/rfc/rfc7540#section-8.3
func RequestURI(r *http.Request) string {
	var uri string
	if r.ProtoMajor == 2 && r.Method == "CONNECT" {
		uri = r.Host
	} else {
		uri = r.RequestURI
	}
	if uri == "" {
		uri = r.URL.RequestURI()
	}
	return uri
}

func (e *Entry) RequestURI() string { return RequestURI(e.Request) }

func Username(r *http.Request) string {
	if r.URL != nil && r.URL.User != nil {
		if user := r.URL.User.Username(); user != "" {
			return user
		}
	}
	return ""
}

func (e *Entry) Username() string { return Username(e.Request) }

func (e *Entry) Timestamp() string {
	return e.Metrics.Time.Format("02/Jan/2006:15:04:05 -0700")
}
