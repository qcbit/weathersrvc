package web

import (
	"context"
	"encoding/json"
	"net/http"
)

// Response converts a Go value to JSON and sends it to the client.
// If the value is of type error, and it is not nil, the error will be
// sent to the client as a JSON object with key "error".
func Response(ctx context.Context, w http.ResponseWriter, data any, statusCode int) error {
	SetStatusCode(ctx, statusCode)

	// If nothing to marshal then return.
	if statusCode == http.StatusNoContent {
		w.WriteHeader(statusCode)
		return nil
	}

	// Convert the response to JSON.
	jsonData, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	// Send the result back to the client.
	if _, err := w.Write(jsonData); err != nil {
		return err
	}

	return nil
}
