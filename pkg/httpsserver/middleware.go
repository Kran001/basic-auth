package httpsserver

import (
	"errors"
	"net/http"
	"time"

	log "github.com/Kran001/basic-auth/pkg/logging"
)

func LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Logger.Infof("Started %s %s", r.Method, r.RequestURI)

		start := time.Now()

		next.ServeHTTP(w, r)

		log.Logger.Infof("Request %s %s finished %v", r.Method, r.RequestURI, time.Since(start))
	})
}

// Recover - Recovering after panic in http handler.
func Recover(h http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		var err error
		defer func() {
			r := recover()
			if r != nil {
				switch t := r.(type) {
				case string:
					err = errors.New(t)
				case error:
					err = t
				default:
					log.Logger.Info("Make recovered error by default")
					err = errors.New("unknown error")
				}
				// we can use it for sending some errors
				log.Logger.Error(err.Error())
				http.Error(writer, err.Error(), http.StatusInternalServerError)
			}
		}()
		h.ServeHTTP(writer, request)
	})
}
