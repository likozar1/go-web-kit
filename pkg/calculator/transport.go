package calculator

import (
	"context"
	"fmt"
	"html/template"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
)

var T *template.Template

func init() {
	t, err := template.ParseFiles("./../../views/index.html", "./../../views/calculator.html")
	if err != nil {
		panic("Template parsing error")
	}
	T = t
}

func decodeRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	expression := r.FormValue("expr")

	return Request{
		Expression: expression,
	}, nil
}

// writes response from endpoint to client
func encodePlainResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	resp := response.(Response)

	w.Header().Add("Content-type", "text/plain")

	_, err := fmt.Fprint(w, resp.Result)
	return err
}

// writes decorated response from endpoint to client
func encodeHTMLResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	resp := response.(Response)

	w.Header().Add("Content-type", "text/html")

	// prepare the data
	data := struct {
		Expr   string
		Result string
	}{
		Expr:   resp.Expression,
		Result: resp.Result,
	}

	return T.ExecuteTemplate(w, "calculator.html", data)
}

// NewHTTPHandler creates greeter handlers
func NewHTTPHandler(endpoint endpoint.Endpoint, logger log.Logger) http.Handler {
	m := http.NewServeMux()

	m.Handle("/api/calculate", httptransport.NewServer(
		EndpointLoggingMiddleware(log.With(logger, "endpoint", "api"))(endpoint),
		decodeRequest,
		encodePlainResponse,
	))

	m.Handle("/calculate", httptransport.NewServer(
		EndpointLoggingMiddleware(log.With(logger, "endpoint", "html"))(endpoint),
		decodeRequest,
		encodeHTMLResponse,
	))

	return m
}
