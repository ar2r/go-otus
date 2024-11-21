package httpserver

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/server"
)

func loggingMiddleware(next http.Handler, logg ServerLogger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rr := &responseRecorder{w, http.StatusOK, 0}
		next.ServeHTTP(rr, r)

		clientIP, _, _ := net.SplitHostPort(r.RemoteAddr)
		if ip, err := server.NormalizeIPv4(clientIP); err == nil {
			clientIP = ip
		}
		timestamp := start.Format("02/Jan/2006:15:04:05 -0700")
		method := r.Method
		path := r.URL.RequestURI()
		protocol := r.Proto
		statusCode := rr.statusCode
		responseSize := rr.responseSize
		userAgent := r.UserAgent()

		logMsg := fmt.Sprintf("%s [%s] %s %s %s %d %d \"%s\"",
			clientIP, timestamp, method, path, protocol, statusCode, responseSize, userAgent)
		logg.Info(logMsg)
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
