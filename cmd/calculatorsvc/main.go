package main

import (
	"flag"
	"net/http"
	"os"

	"github.com/go-kit/kit/log"

	"github.com/go-web-kit/pkg/calculator"
)

func main() {
	var addr = flag.String("addr", "127.0.0.1:8080", "Interface and port to listen on")

	// parse the flags
	flag.Parse()

	// Create a single logger, which we'll use and give to other components.
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	// Create calculator service
	var service calculator.Service
	{
		service = calculator.New()
		service = calculator.ServiceLoggingMiddleware(log.With(logger, "service", "calculator"))(service)
		service = calculator.ValidateMiddleware()(service)
	}

	endpoint := calculator.MakeEndpoint(service)

	handler := calculator.NewHTTPHandler(endpoint, logger)

	logger.Log("transport", "http", "listen", *addr)
	err := http.ListenAndServe(*addr, handler)
	if err != nil {
		logger.Log("transport", "http", "during", "listen", "err", err)
	}
}
