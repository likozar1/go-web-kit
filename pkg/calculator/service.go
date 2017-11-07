package calculator

import (
	"regexp"
	"strconv"
)

const Regex = "-?[0-9]+|[*+/-]"

var Re *regexp.Regexp

func init() {
	Re = regexp.MustCompile(Regex)
}

type calculator string

// Service parses and calculates expression
type Service interface {
	ParseExpression(expression string) []string
	CalculateParsedExpression(expressionArray []string) string
}

// New service
func New() Service {
	return calculator("expression")
}

func (c calculator) ParseExpression(expression string) []string {
	parsed := Re.FindAllString(expression, -1)

	if len(parsed)%2 != 1 {
		panic("invalid expression")
	}

	return parsed
}

func calculate(first int, second int, operator string) int {
	result := 0
	switch {
	case operator == "+":
		result = first + second
	case operator == "-":
		result = first - second
	case operator == "*":
		result = first * second
	case operator == "/":
		if second == 0 {
			panic("Division by zero")
		}
		result = first / second
	default:
		panic("undefined operator")
	}

	return result
}

func (c calculator) CalculateParsedExpression(expressionArray []string) string {
	for i := 1; i < len(expressionArray); i = i + 2 {
		first, err := strconv.Atoi(expressionArray[i-1])
		second, err := strconv.Atoi(expressionArray[i+1])

		if err != nil {
			panic("invalid expression")
		}

		expressionArray[i+1] = strconv.Itoa(calculate(first, second, expressionArray[i]))
		expressionArray[i-1] = ""
		expressionArray[i] = ""
	}

	return expressionArray[len(expressionArray)-1]
}
