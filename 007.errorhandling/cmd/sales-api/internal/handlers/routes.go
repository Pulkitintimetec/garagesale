package handlers

import (
	"log"
	"net/http"
	"os"

	"garagesale/007.errorhandling/internal/mid"
	"garagesale/007.errorhandling/internal/platform/auth"
	"garagesale/007.errorhandling/internal/platform/web"
)

// API construts a handler that knows about all api routes
func API(shutdowns chan os.Signal, logger *log.Logger, authenticator *auth.Authenticator) http.Handler {
	app := web.NewApp(shutdowns, logger, mid.Logger(logger), mid.Errors(logger), mid.Metrics(logger), mid.Panics(logger))
	{
		h := Heartbeat{}
		app.Handle(http.MethodGet, "/heart", h.Health)
	}

	{
		// Register user handlers.
		u := Users{
			authenticator: authenticator,
		}
		app.Handle(http.MethodGet, "/users/token", u.Token)
	}
	{
		p := Product{
			Log: logger,
		}

		app.Handle(http.MethodGet, "/getProducts", p.List, mid.Authenticate(authenticator))
		app.Handle(http.MethodGet, "/getProductByName/{name}/{cost}", p.Retrieve, mid.Authenticate(authenticator))
		app.Handle(http.MethodGet, "/ /{productID}", p.GetProductUsingID, mid.Authenticate(authenticator))
		app.Handle(http.MethodPut, "/updateProduct/{productID}", p.Update, mid.Authenticate(authenticator), mid.HasRole(auth.RoleAdmin))
		app.Handle(http.MethodPost, "/insertProduct", p.Insert, mid.Authenticate(authenticator))
		app.Handle(http.MethodPost, "/addSale/{productID}/sale", p.AddSale, mid.Authenticate(authenticator))
		app.Handle(http.MethodGet, "/listSale/{productID}", p.ListSale, mid.Authenticate(authenticator))
		app.Handle(http.MethodDelete, "/deleteProduct/{productID}", p.DeleteProductByID, mid.Authenticate(authenticator), mid.HasRole(auth.RoleAdmin))

	}
	return app
}
