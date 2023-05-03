package httputil

import (
	"encoding/json"
	"net/http"

	"github.com/jseow5177/pockteer-be/pkg/errutil"
)

type Response struct {
	Code  int         `json:"code"`
	Error string      `json:"error"`
	Body  interface{} `json:"body"`
}

func ReturnServerResponse(w http.ResponseWriter, res interface{}, resErr error) {
	code, errMsg := errutil.ParseHttpError(resErr)

	resp := &Response{
		Code:  code,
		Error: errMsg,
		Body:  res,
	}

	js, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(js)
}
