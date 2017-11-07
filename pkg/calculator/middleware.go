package calculator

import (
	"github.com/go-kit/kit/log"
)

// Middleware describes a service (as opposed to endpoint) middleware.
type Middleware func(Service) Service

// ValidateMiddleware is a validator middleware for service
func ValidateMiddleware() Middleware {
	return func(next Service) Service {
		return validateMiddleware{next}
	}
}

type validateMiddleware struct {
	next Service
}

func (mw validateMiddleware) CalculateParsedExpression(expression []string) string {
	if len(expression) == 0 {
		return ""
	}

	return mw.next.CalculateParsedExpression(expression)
}

func (mw validateMiddleware) ParseExpression(expression string) []string {
	if expression == "" {
		return []string{}
	}
	return mw.next.ParseExpression(expression)
}

// ServiceLoggingMiddleware is a validator middleware for service
func ServiceLoggingMiddleware(log log.Logger) Middleware {
	return func(next Service) Service {
		return loggingMiddleware{log, next}
	}
}

type loggingMiddleware struct {
	logger log.Logger
	next   Service
}

func (mw loggingMiddleware) CalculateParsedExpression(expression []string) string {
	mw.logger.Log("method", "CalculateParsedExpression", "expression", expression)
	return mw.next.CalculateParsedExpression(expression)
}

func (mw loggingMiddleware) ParseExpression(expression string) []string {
	mw.logger.Log("method", "ParseExpression", "expression", expression)
	return mw.next.ParseExpression(expression)
}
