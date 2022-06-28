package internalhttp

import (
	"net/http"
	"strconv"
	"time"
)

func loggingMiddleware(h http.HandlerFunc, s *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.logger.Info(r.RemoteAddr + "[" + time.Now().UTC().String() + "]" + r.Method + r.RequestURI + strconv.Itoa(int(r.ContentLength)) + r.UserAgent())
		h(w, r)
	}
}
