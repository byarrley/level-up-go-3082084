package main

import (
	"log"
	"strings"
	"time"
)

const delay = 700 * time.Millisecond

// print outputs a message and then sleeps for a pre-determined amount
func print(msg string) {
	log.Println(msg)
	time.Sleep(delay)
}

// slowDown takes the given string and repeats its characters
// according to their index in the string.
func slowDown(msg string) {
	sl := strings.Split(msg, " ")

	for _, s := range sl {
		var m []string
		for i, b := range s {
			//Solution: this loop can be eliminated by using strings.Repeat(.)
			for idx := 0; idx <= i; idx++ {
				m = append(m, string(b))
			}
		}
		print(strings.Join(m, ""))
	}

	//panic("NOT IMPLEMENTED")
}

func main() {
	msg := "Time to learn about Go strings!"
	slowDown(msg)
}
