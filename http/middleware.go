package astihttp

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/asticode/go-astilog"
	"github.com/pkg/errors"
)

// ChainMiddlewares chains middlewares
func ChainMiddlewares(h http.Handler, ms ...Middleware) http.Handler {
	for _, m := range ms {
		h = m(h)
	}
	return h
}

// Middleware represents a middleware
type Middleware func(http.Handler) http.Handler

// MiddlewareTimeout adds a timeout to a handler
func MiddlewareTimeout(timeout time.Duration) Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			// Init context
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()

			// Serve
			var done = make(chan bool)
			go func() {
				h.ServeHTTP(rw, r)
				done <- true
			}()

			// Wait for done or timeout
			for {
				select {
				case <-ctx.Done():
					astilog.Error(errors.Wrap(ctx.Err(), "serving HTTP failed"))
					rw.WriteHeader(http.StatusGatewayTimeout)
					return
				case <-done:
					return
				}
			}
		})
	}
}

// MiddlewareBasicAuth adds basic HTTP auth to a handler
func MiddlewareBasicAuth(username, password, prefix string) Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			// Only authenticate if prefix is correct
			if prefix == "" || strings.HasPrefix(r.URL.EscapedPath(), prefix) {
				if u, p, ok := r.BasicAuth(); !ok || u != username || p != password {
					rw.Header().Set("WWW-Authenticate", "Basic Realm")
					rw.WriteHeader(http.StatusUnauthorized)
					return
				}
			}

			// Next handler
			h.ServeHTTP(rw, r)
		})
	}
}
