package router

import (
	"context"
	"net/http"
	"reflect"

	"github.com/gorilla/mux"

	"github.com/jseow5177/pockteer-be/pkg/errutil"
	"github.com/jseow5177/pockteer-be/pkg/httputil"
	"github.com/jseow5177/pockteer-be/pkg/validator"
)

type Middleware interface {
	Handle(http.Handler) http.Handler
}

type Handler struct {
	Req        interface{}
	Res        interface{}
	Validator  validator.Validator
	HandleFunc func(ctx context.Context, req interface{}, res interface{}) error

	reqT  reflect.Type
	respT reflect.Type
}

type HttpRoute struct {
	Method      string
	Path        string
	Handler     Handler
	Middlewares []Middleware
}

type HttpRouter struct {
	*mux.Router
}

func (r *HttpRouter) RegisterHttpRoute(hr *HttpRoute) {
	// save req and res type
	hr.Handler.reqT = reflect.TypeOf(hr.Handler.Req).Elem()
	hr.Handler.respT = reflect.TypeOf(hr.Handler.Res).Elem()

	// calling chain
	chain := http.Handler(hr.Handler)

	if hr.Middlewares != nil {
		// wrap middlewares from right to left
		for i := len(hr.Middlewares) - 1; i >= 0; i-- {
			chain = hr.Middlewares[i].Handle(chain)
		}
	}

	r.Methods(hr.Method).Path(hr.Path).Handler(chain)
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	req := reflect.New(h.reqT).Interface()
	res := reflect.New(h.respT).Interface()

	// parse body
	if r.Body != http.NoBody {
		err := httputil.ReadJsonBody(r, req)
		if err != nil {
			httputil.ReturnServerResponse(w, nil, errutil.BadRequestError(err))
			return
		}
	}

	if h.Validator != nil {
		err := h.Validator.Validate(req)
		if err != nil {
			httputil.ReturnServerResponse(w, nil, errutil.ValidationError(err))
			return
		}
	}

	err := h.HandleFunc(r.Context(), req, res)
	httputil.ReturnServerResponse(w, res, err)
}
