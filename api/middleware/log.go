package middleware

import (
	"net/http"
	"time"

	"github.com/rs/zerolog/hlog"
	"github.com/rs/zerolog/log"
)

func Log(handler http.Handler) http.Handler {
	// add request ID
	logReqID := hlog.RequestIDHandler("req_id", "Request-ID")

	// calls f after a request is done
	f := func(r *http.Request, status, _ int, duration time.Duration) {
		hlog.FromRequest(r).
			Info().
			Str("method", r.Method).
			Stringer("url", r.URL).
			Dur("duration", duration).
			Int("status", status).
			Msg("")
	}

	return hlog.NewHandler(log.Logger)(logReqID(hlog.AccessHandler(f)(handler)))
}
