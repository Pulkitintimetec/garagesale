package product_test

import (
	"context"
	"testing"
	"time"

	"garagesale/007.errorhandling/internal/platform/auth"
	"garagesale/007.errorhandling/internal/product"
	"github.com/google/go-cmp/cmp"
)

//TestProducts ...
func TestProducts(t *testing.T) {
	var UserID auth.Claims
	newP := product.NewProduct{
		Name:     "Comic Book",
		Cost:     "10",
		Quantity: "55",
	}
	ctx := context.Background()
	now := time.Date(2020, time.March, 1, 0, 0, 0, 0, time.UTC)
	p0, err := product.InsertProduct(ctx, UserID, newP, now)
	if err != nil {
		t.Fatalf("creating product p0: %s", err)
	}

	p1, err := product.GetDatabyName(ctx, p0.Name, p0.Cost)
	if err != nil {
		t.Fatalf("getting product p0: %s", err)
	}

	if diff := cmp.Diff(p1, p0); diff != "" {
		t.Fatalf("fetched != created:\n%s", diff)
	}
}

func TestProductList(t *testing.T) {
	ctx := context.Background()
	ps, err := product.GetAllData(ctx)
	if err != nil {
		t.Fatalf("listing products: %s", err)
	}
	if exp, got := 7, len(ps); exp != got {
		t.Fatalf("expected product list size %v, got %v", exp, got)
	}
}
