package web

import (
	"context"
	"log"
	"net/http"
	"os"
	"syscall"
	"time"

	"github.com/go-chi/chi"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/plugin/ochttp/propagation/tracecontext"
	"go.opencensus.io/trace"
)

type ctxKey int

const KeyValues ctxKey = 1

// App is the entrypoint into our application and what controls the context of
// each request. Feel free to add any configuration data/logic on this type.
type App struct {
	log      *log.Logger
	mux      *chi.Mux
	mw       []Middleware
	shutdown chan os.Signal
	och      *ochttp.Handler
}

// Values carries inforamation of each request
type Values struct {
	TraceID    string
	StatusCode int
	Start      time.Time
}

// Handler is the signature that all handler will imlement
type Handler func(context.Context, http.ResponseWriter, *http.Request) error

// NewApp constructs an App to handle a set of routes.
func NewApp(shutdown chan os.Signal, log *log.Logger, mw ...Middleware) *App {
	app := App{
		log:      log,
		mux:      chi.NewRouter(),
		mw:       mw,
		shutdown: shutdown,
	}

	// Create an OpenCensus HTTP Handler which wraps the router. This will start
	// the initial span and annotate it with information about the request/response.
	//
	// This is configured to use the W3C TraceContext standard to set the remote
	// parent if an client request includes the appropriate headers.
	// https://w3c.github.io/trace-context/
	app.och = &ochttp.Handler{
		Handler:     app.mux,
		Propagation: &tracecontext.HTTPFormat{},
	}

	return &app
}

// Handle associates a handler function with an HTTP Method and URL pattern.
func (a *App) Handle(method, url string, h Handler, mw ...Middleware) {

	// First wrap handler specific middleware around this handler.
	h = wrapMiddleware(mw, h)

	// Add the application's general middleware to the handler chain.
	h = wrapMiddleware(a.mw, h)

	fn := func(w http.ResponseWriter, r *http.Request) {
		v := Values{
			Start: time.Now(),
		}
		ctx := context.WithValue(r.Context(), KeyValues, &v)
		ctx, span := trace.StartSpan(ctx, "internal.platform.web")
		defer span.End()
		if err := h(ctx, w, r); err != nil {
			a.log.Printf("%s : unhandled error: %+v", v.TraceID, err)
			if IsShutdown(err) {
				a.SignalShutdown()
			}
		}
	}
	a.mux.MethodFunc(method, url, fn)
}

// ServeHTTP implements the http.Handler interface.
func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.och.ServeHTTP(w, r)
}

// SignalShutdown is used to gracefully shutdown the app when an integrity
// issue is identified.
func (a *App) SignalShutdown() {
	a.log.Println("error returned from handler indicated integrity issue, shutting down service")
	a.shutdown <- syscall.SIGHUP
}
