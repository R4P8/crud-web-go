package middleware

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// Wrap converts a httprouter.Handle to http.Handler with otel tracing support
func Wrap(h httprouter.Handle, operationName string) http.Handler {
	return otelhttp.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := httprouter.ParamsFromContext(r.Context())
		h(w, r, params)
	}), operationName)
}

