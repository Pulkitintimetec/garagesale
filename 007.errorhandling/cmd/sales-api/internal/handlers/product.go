package handlers

import (
	"context"
	"log"
	"net/http"
	"time"

	"garagesale/007.errorhandling/internal/platform/auth"
	"garagesale/007.errorhandling/internal/platform/web"
	"garagesale/007.errorhandling/internal/product"
	"github.com/go-chi/chi"
	"github.com/pkg/errors"
	"go.opencensus.io/trace"
)

// Product ...
type Product struct {
	Log *log.Logger
	Ctx context.Context
}

// List give all products from the database as a list
func (p *Product) List(ctx context.Context, w http.ResponseWriter, req *http.Request) error {
	ctx, span := trace.StartSpan(ctx, "handlers.product.List")
	defer span.End()
	// var br BookRepo
	p.Log.Println("Testing")
	getdata, err := product.GetAllData(req.Context())
	if err != nil {
		return err
	}
	return web.Respond(ctx, w, getdata, http.StatusOK)

}

// Retrieve give particular products from the database
func (p *Product) Retrieve(ctx context.Context, w http.ResponseWriter, req *http.Request) error {
	// var br BookRepo
	p.Log.Println("Testing")
	Name := chi.URLParam(req, "name")
	Cost := chi.URLParam(req, "cost")
	retrieveData, err := product.GetDatabyName(req.Context(), Name, Cost)
	if err != nil {
		return err
	}
	return web.Respond(ctx, w, retrieveData, http.StatusOK)

}

//Update ...
func (p *Product) Update(ctx context.Context, w http.ResponseWriter, req *http.Request) error {
	id := chi.URLParam(req, "productID")

	var update product.UpdateProductStructure
	if err := web.Decode(req, &update); err != nil {
		return errors.Wrap(err, "decoding product update")
	}
	claims, ok := ctx.Value(auth.Key).(auth.Claims)
	if !ok {
		return errors.New("claims missing from context")
	}

	if err := product.UpdateProduct(req.Context(), claims, id, update, time.Now()); err != nil {
		switch err {
		case product.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case product.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		case product.ErrForbidden:
			return web.NewRequestError(err, http.StatusForbidden)
		default:
			return errors.Wrapf(err, "updating product %q", id)
		}
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)

}

// GetProductUsingID is the handler function of getting product with passing ID
func (p *Product) GetProductUsingID(ctx context.Context, w http.ResponseWriter, req *http.Request) error {
	// var br BookRepo
	productID := chi.URLParam(req, "productID")
	p.Log.Println("Listing Of Product Particular to ID")
	getdata, err := product.GetProductByID(req.Context(), productID)
	if err != nil {
		return err
	}
	return web.Respond(ctx, w, getdata, http.StatusOK)

}

// DeleteProductByID deleting Handler
func (p *Product) DeleteProductByID(ctx context.Context, w http.ResponseWriter, req *http.Request) error {
	productID := chi.URLParam(req, "productID")
	p.Log.Println("Deleting the Product Particular to ID")
	err := product.DeleteProduct(req.Context(), productID)
	if err != nil {
		return err
	}
	return web.Respond(ctx, w, nil, http.StatusOK)
}

// Insert give particular products from the database
func (p *Product) Insert(ctx context.Context, w http.ResponseWriter, req *http.Request) error {
	// var br BookRepo
	p.Log.Println("Testing")
	var np product.NewProduct
	var now time.Time
	claims, ok := ctx.Value(auth.Key).(auth.Claims)
	if !ok {
		return errors.New("auth claims are not in context")
	}
	if err := web.Decode(req, &np); err != nil {
		return err
	}
	insertData, err := product.InsertProduct(req.Context(), claims, np, now)
	if err != nil {
		return err
	}
	return web.Respond(ctx, w, insertData, http.StatusOK)
}

// AddSale use to add sale in the database with respective to productId
func (p *Product) AddSale(ctx context.Context, w http.ResponseWriter, req *http.Request) error {
	// var br BookRepo
	p.Log.Println("AddingSale")
	var ns product.NewSale
	var now time.Time
	productID := chi.URLParam(req, "productID")
	if err := web.Decode(req, &ns); err != nil {
		return err
	}

	insertData, err := product.AddSales(req.Context(), ns, productID, now)
	if err != nil {
		return err
	}
	return web.Respond(ctx, w, insertData, http.StatusOK)
}

// ListSale used to give all sales
func (p *Product) ListSale(ctx context.Context, w http.ResponseWriter, req *http.Request) error {
	// var br BookRepo
	productID := chi.URLParam(req, "productID")
	p.Log.Println("Listing Of Sale")
	getdata, err := product.ListS(req.Context(), productID)
	if err != nil {
		return err
	}
	return web.Respond(ctx, w, getdata, http.StatusOK)

}
