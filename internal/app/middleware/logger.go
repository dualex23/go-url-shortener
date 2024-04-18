package middleware

import (
	"net/http"
	"time"

	"github.com/dualex23/go-url-shortener/internal/app/utils"
)

type responseData struct {
	status int
	size   int
}
type loggingResponseWriter struct {
	http.ResponseWriter
	responseData *responseData
}

func (l *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := l.ResponseWriter.Write(b)
	l.responseData.size += size
	return size, err
}
func (l *loggingResponseWriter) WriteHeader(statuscode int) {
	l.ResponseWriter.WriteHeader(statuscode)
	l.responseData.status = statuscode
}

func WithLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		start := time.Now()

		response := &responseData{
			status: 0,
			size:   0,
		}

		lw := loggingResponseWriter{
			ResponseWriter: w,
			responseData:   response,
		}

		uri := r.RequestURI
		method := r.Method

		next.ServeHTTP(&lw, r)

		duration := time.Since(start)

		utils.GetLogger().Infoln(
			"uri", uri,
			"method", method,
			"duration", duration,
			"statuscode", response.status,
			"size", response.size,
		)
	})
}
