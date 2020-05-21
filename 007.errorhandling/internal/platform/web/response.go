package web

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

// Respond marshal a values to JSON sends it to client
func Respond(ctx context.Context, w http.ResponseWriter, val interface{}, statusCode int) error {
	v, ok := ctx.Value(KeyValues).(*Values)
	if !ok {
		return errors.New("web values missing from context")
	}
	v.StatusCode = statusCode
	data, err := json.Marshal(val)
	if err != nil {
		return errors.Wrap(err, "marshalling value to JSON")
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode) // 200
	if _, err := w.Write(data); err != nil {

		return errors.Wrap(err, "writing to client")
	}
	return nil
}

//ResponseError knows how to handle errors going out to the client
func ResponseError(ctx context.Context, w http.ResponseWriter, err error) error {

	// If the error was of the type *Error, the handler has
	// a specific status code and error to return.
	if webErr, ok := errors.Cause(err).(*Error); ok {
		er := ErrorResponse{
			Error:  webErr.Err.Error(),
			Fields: webErr.Fields,
		}
		if err := Respond(ctx, w, er, webErr.Status); err != nil {
			return err
		}
		return nil
	}

	// If not, the handler sent any arbitrary error value so use 500.
	er := ErrorResponse{
		Error: http.StatusText(http.StatusInternalServerError),
	}
	if err := Respond(ctx, w, er, http.StatusInternalServerError); err != nil {
		return err
	}
	return nil
}
