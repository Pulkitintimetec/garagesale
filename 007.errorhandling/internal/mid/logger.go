package mid

import (
	"context"
	"log"
	"net/http"
	"time"

	"garagesale/007.errorhandling/internal/platform/web"
	"go.opencensus.io/trace"
)

// Logger writes some information about the request to the logs in the
// format: (200) GET /foo -> IP ADDR (latency)
func Logger(log *log.Logger) web.Middleware {

	// This is the actual middleware function to be executed.
	f := func(before web.Handler) web.Handler {

		// Create the handler that will be attached in the middleware chain.
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			ctx, span := trace.StartSpan(ctx, "internal.mid.logger")
			defer span.End()
			v, ok := ctx.Value(web.KeyValues).(*web.Values)
			if !ok {
				return web.NewShutdownError("web values missing from context")
			}
			err := before(ctx, w, r)

			log.Printf("%d %s %s (%v) ", v.StatusCode, r.Method, r.URL.Path, time.Since(v.Start))

			// Return the error so it can be handled further up the chain.
			return err
		}

		return h
	}

	return f
}
