package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	httpModels "github.com/MoonSHRD/shortify/models/http"
)

func ReturnHTTPError(w http.ResponseWriter, err string, code int) {
	res := httpModels.HTTPError{
		Error: httpModels.ErrorMessage{
			Message: err,
		},
	}
	r, _ := json.Marshal(res)
	writeHTTPJSONError(w, string(r), code)
}

func writeHTTPJSONError(w http.ResponseWriter, err string, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	fmt.Fprintln(w, err)
}
