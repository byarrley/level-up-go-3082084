package main

import (
	"flag"
	"log"
	"slices"
)

// isBalanced returns whether the given expression
// has balanced brackets.

// Task: "Given a string mathematical expression, implement a function that outputs if the expression has balanced brackets".  Apparently I didn't understand the problem, so I will update this solution.
// balanced: "1 + (2 + 3)"
// unbalanced: "[1 + ( 3 * 2 ])"
type opType int

const (
	openOp opType = iota
	closeOp
	otherOp
)

var validPairs = map[rune]rune{
	'(': ')',
	'[': ']',
	'{': '}',
}

func getOpType(op rune) opType {
	for o, c := range validPairs {
		switch op {
		case o:
			return openOp
		case c:
			return closeOp
		}
	}
	return otherOp
}

func isBalanced(expr string) bool {
	stack := Stack{}

	for _, r := range expr {
		switch getOpType(r) {
		case openOp:
			stack.Push(r)
		case closeOp:
			tmp := stack.Pop()
			if validPairs[tmp] != r {
				return false
			}
		default:
			continue
		}
	}
	return stack.IsEmpty()
}

// printResult prints whether the expression is balanced.
func printResult(expr string, balanced bool) {
	if balanced {
		log.Printf("%s is balanced.\n", expr)
		return
	}
	log.Printf("%s is not balanced.\n", expr)
}

func main() {
	expr := flag.String("expr", "", "The expression to validate brackets on.")
	flag.Parse()
	printResult(*expr, isBalanced(*expr))
}

type Stack struct {
	data []rune
	size int
}

func (s *Stack) Push(i rune) {
	s.data = append(s.data, i)
	s.size++
}

func (s *Stack) Pop() rune {
	if s.IsEmpty() {
		return 0
	}

	tmp := s.data[s.size-1]
	s.data = slices.Delete(s.data, s.size-1, s.size)
	s.size--
	return tmp
}

func (s *Stack) Peek() rune {
	return s.data[s.size-1]
}

func (s *Stack) IsEmpty() bool {
	return s.Size() == 0
}

func (s *Stack) Size() int {
	return s.size
}
