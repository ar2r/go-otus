package internalhttp

import (
	"fmt"
	"net/http"
	"time"
)

func loggingMiddleware(next http.Handler, logg Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rr := &responseRecorder{w, http.StatusOK, 0}
		next.ServeHTTP(rr, r)

		clientIP := r.RemoteAddr
		timestamp := start.Format("01/Jan/2000:23:59:59 +0300")
		method := r.Method
		path := r.URL.RequestURI()
		protocol := r.Proto
		statusCode := rr.statusCode
		responseSize := rr.responseSize
		userAgent := r.UserAgent()

		logMsg := fmt.Sprintf("%s [%s] %s %s %s %d %d \"%s\"",
			clientIP, timestamp, method, path, protocol, statusCode, responseSize, userAgent)
		logg.InfoRaw(logMsg)
	})
}

type responseRecorder struct {
	http.ResponseWriter
	statusCode   int
	responseSize int
}

func (rr *responseRecorder) WriteHeader(statusCode int) {
	rr.statusCode = statusCode
	rr.ResponseWriter.WriteHeader(statusCode)
}

func (rr *responseRecorder) Write(b []byte) (int, error) {
	size, err := rr.ResponseWriter.Write(b)
	rr.responseSize += size
	return size, err
}
