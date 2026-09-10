package main

import (
	"flag"
	"log"
	"strconv"
	"strings"
)

// The Task: Given a simple calculator, implement edge case and input verification that outputs error messages and gracefully shuts down.
/* Plan:
* Features:
	- Calculator only supports the 4 basic operators +,-,*,/
	- Calculator only supports two numerical operands at a time
* Identify possible edge cases
	- Invalid operands
	- Invalid operators
	- Not enough tokens in the expression
	- Too many tokens in the expression
	- Tokens not separated by spaces
*/

// operators is the map of legal operators and their functions
var operators = map[string]func(x, y float64) float64{
	"+": func(x, y float64) float64 { return x + y },
	"-": func(x, y float64) float64 { return x - y },
	"*": func(x, y float64) float64 { return x * y },
	"/": func(x, y float64) float64 { return x / y },
}

// parseOperand parses a string to a float64
func parseOperand(op string) float64 {
	parsedOp, _ := strconv.ParseFloat(op, 64)
	return parsedOp
}

// calculate returns the result of a 2 operand mathematical expression
func calculate(expr string) float64 {
	ops := strings.Fields(expr)
	left := parseOperand(ops[0])
	right := parseOperand(ops[2])
	f := operators[ops[1]]
	result := f(left, right)
	return result
}

func main() {
	expr := flag.String("expr", "",
		"The expression to calculate on, separated by spaces.")
	flag.Parse()
	result := calculate(*expr)
	log.Printf("%s = %.2f\n", *expr, result)
}
