package tests

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	// NOTE: Models should not be imported, we want to test the exact JSON. We
	// make the comparison process easier using the go-cmp library.
	// "github.com/ardanlabs/garagesale/cmd/sales-api/internal/handlers"
	"garagesale/007.errorhandling/cmd/sales-api/internal/handlers"
	"garagesale/007.errorhandling/internal/platform/auth"
	"github.com/google/go-cmp/cmp"
)

// TestProducts runs a series of tests to exercise Product behavior from the
// API level. The subtests all share the same database and application for
// speed and convenience. The downside is the order the tests are ran matters
// and one test may break if other tests are not ran before it. If a particular
// subtest needs a fresh instance of the application it can make it or it
// should be its own Test* function.
func TestProducts(t *testing.T) {

	log := log.New(os.Stderr, "TEST : ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)
	var authenticator *auth.Authenticator
	shutdown := make(chan os.Signal, 1)
	tests := ProductTests{app: handlers.API(shutdown, log, authenticator)}

	//t.Run("List", tests.List)
	t.Run("ProductCRUD", tests.ProductCRUD)
}

// ProductTests holds methods for each product subtest. This type allows
// passing dependencies for tests while still providing a convenient syntax
// when subtests are registered.
type ProductTests struct {
	app http.Handler
}

func (p *ProductTests) List(t *testing.T) {
	req := httptest.NewRequest("GET", "/getProducts", nil)
	resp := httptest.NewRecorder()

	p.app.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("getting: expected status code %v, got %v", http.StatusOK, resp.Code)
	}

	var list []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&list); err != nil {
		t.Fatalf("decoding: %s", err)
	}

	want := []map[string]interface{}{
		{
			"Cost":         "15000",
			"Name":         "Samsung M30s",
			"Quantity":     "50",
			"date_created": "0001-01-01T00:00:00Z",
			"date_updated": "0001-01-01T00:00:00Z",
		},
		{
			"Cost":         "55000",
			"Name":         "Sofa Set",
			"Quantity":     "5",
			"date_created": "0001-01-01T00:00:00Z",
			"date_updated": "0001-01-01T00:00:00Z",
		},
		{
			"Cost":         "19000",
			"Name":         "MicroWave",
			"Quantity":     "10",
			"date_created": "0001-01-01T00:00:00Z",
			"date_updated": "0001-01-01T00:00:00Z",
		},
		{
			"Cost":         "10",
			"Name":         "Comic Book",
			"Quantity":     "55",
			"date_created": "0001-01-01T00:00:00Z",
			"date_updated": "0001-01-01T00:00:00Z",
		},
		{
			"Cost":         "1000",
			"Name":         "Helmet",
			"Quantity":     "35",
			"date_created": "0001-01-01T00:00:00Z",
			"date_updated": "0001-01-01T00:00:00Z",
		},
		{
			"Cost":         "55",
			"Name":         "product0",
			"Quantity":     "6",
			"date_created": "0001-01-01T00:00:00Z",
			"date_updated": "0001-01-01T00:00:00Z",
		},
	}

	if diff := cmp.Diff(want, list); diff != "" {
		t.Fatalf("Response did not match expected. Diff:\n%s", diff)
	}
}

func (p *ProductTests) ProductCRUD(t *testing.T) {
	var created map[string]interface{}

	{ // CREATE
		body := strings.NewReader(`{"Name":"product01","Cost":"55","Quantity":"6"}`)

		req := httptest.NewRequest("POST", "/insertProduct", body)
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		p.app.ServeHTTP(resp, req)

		if http.StatusOK != resp.Code {
			t.Fatalf("posting: expected status code %v, got %v", http.StatusCreated, resp.Code)
		}

		if err := json.NewDecoder(resp.Body).Decode(&created); err != nil {
			t.Fatalf("decoding: %s", err)
		}
		if created["date_created"] == "" || created["date_created"] == nil {
			t.Fatal("expected non-empty product date_created")
		}
		if created["date_updated"] == "" || created["date_updated"] == nil {
			t.Fatal("expected non-empty product date_updated")
		}

		want := map[string]interface{}{

			"Name":         "product01",
			"Cost":         "55",
			"Quantity":     "6",
			"date_created": created["date_created"],
			"date_updated": created["date_updated"],
		}

		if diff := cmp.Diff(want, created); diff != "" {
			t.Fatalf("Response did not match expected. Diff:\n%s", diff)
		}
	}

	{ // READ
		url := fmt.Sprintf("/getProductByName/%s/%s", created["Name"], created["Cost"])
		req := httptest.NewRequest("GET", url, nil)
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		p.app.ServeHTTP(resp, req)

		if http.StatusOK != resp.Code {
			t.Fatalf("retrieving: expected status code %v, got %v", http.StatusOK, resp.Code)
		}

		var fetched map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&fetched); err != nil {
			t.Fatalf("decoding: %s", err)
		}

		// Fetched product should match the one we created.
		if diff := cmp.Diff(created, fetched); diff != "" {
			t.Fatalf("Retrieved product should match created. Diff:\n%s", diff)
		}
	}
}
