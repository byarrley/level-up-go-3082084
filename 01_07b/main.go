package main

import (
	"flag"
	"fmt"
	"log"
	"slices"
)

// isBalanced returns whether the given expression
// has balanced brackets.
func isBalanced(expr string) bool {
	var round, square, curly Stack

	//Solution video: I like what the instructor did with the enumerated types (I was thinking something similar), but the logic doesn't work if the closing bracket doesn't follow the opening bracket
	// As in this case: "(blah (( {} [[ )) ]] "...which I guess isn't terribly useful.
	//Task: "Given a string mathematical expression, implement a function that outputs if the expression has balanced brackets".  Apparently I didn't understand the problem, so I will update this solution.
	for _, r := range expr {
		switch r {
		case '(':
			round.Push(r)
		case '[':
			square.Push(r)
		case '{':
			curly.Push(r)
		case ')':
			round.Pop(r)
		case ']':
			square.Pop(r)
		case '}':
			curly.Pop(r)
		}
	}
	return round.IsEmpty() && square.IsEmpty() && curly.IsEmpty()
	//panic("NOT IMPLEMENTED")
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

func (s *Stack) Pop(i rune) (rune, error) {
	if s.IsEmpty() {
		err := fmt.Errorf("Cannot pop from empty stack!")
		return 0, err
	}

	tmp := s.data[s.size-1]
	s.data = slices.Delete(s.data, s.size-1, s.size)
	s.size--
	return tmp, nil
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
