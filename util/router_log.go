package util

import (
	"github.com/go-chi/chi/middleware"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"net"
	"net/http"
	"time"
)

// Logger returns a request logging middleware
func Logger(category string) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			reqID := middleware.GetReqID(r.Context())
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			t1 := time.Now()
			defer func() {
				remoteIP, _, err := net.SplitHostPort(r.RemoteAddr)
				if err != nil {
					remoteIP = r.RemoteAddr
				}
				scheme := "http"
				if r.TLS != nil {
					scheme = "https"
				}
				fields := logrus.Fields{
					"status_code":      ww.Status(),
					"bytes":            ww.BytesWritten(),
					"duration":         int64(time.Since(t1)),
					"duration_display": time.Since(t1).String(),
					"category":         category,
					"remote_ip":        remoteIP,
					"proto":            r.Proto,
					"method":           r.Method,
				}
				if len(reqID) > 0 {
					fields["request_id"] = reqID
				}
				log.WithFields(fields).Infof("%s://%s%s", scheme, r.Host, r.RequestURI)
			}()

			h.ServeHTTP(ww, r)
		}

		return http.HandlerFunc(fn)
	}
}
