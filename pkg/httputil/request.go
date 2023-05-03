package httputil

import (
	"encoding/json"
	"net/http"
)

func ReadJsonBody(r *http.Request, dst interface{}) error {
	d := json.NewDecoder(r.Body)

	return d.Decode(dst)
}
