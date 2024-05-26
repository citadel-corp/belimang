package middleware

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"github.com/citadel-corp/go-project-template/internal/common/response"
	"github.com/rs/zerolog/log"
)

type LogResponseWriter struct {
	http.ResponseWriter
	statusCode int
	buf        bytes.Buffer
}

func NewLogResponseWriter(w http.ResponseWriter) *LogResponseWriter {
	return &LogResponseWriter{ResponseWriter: w}
}

func (w *LogResponseWriter) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *LogResponseWriter) Write(body []byte) (int, error) {
	w.buf.Write(body)
	return w.ResponseWriter.Write(body)
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		// requestID := id.GenerateStringID(12)
		logRespWriter := NewLogResponseWriter(w)
		next.ServeHTTP(logRespWriter, r)

		log.Debug().
			Dur("duration", time.Since(startTime)).
			Int("status", logRespWriter.statusCode).
			Str("uri", r.RequestURI).
			Str("method", r.Method).
			Msg("request information")
		var resp response.ResponseBody
		err := json.NewDecoder(&logRespWriter.buf).Decode(&resp)
		if logRespWriter.statusCode >= 500 && err == nil {
			log.Error().
				Str("error", resp.Error).
				Msg("internal server errror on request")
		}
	})
}
