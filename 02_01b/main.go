package main

import (
	"flag"
	"log"
	"sync"
)

// The Task: Given a list of messages and a number N, implement a function that outputs the same message N times concurrently
// Hint: Goroutines and channels are the Go concurrency mechanisms
var messages = []string{
	"Hello!",
	"How are you?",
	"Are you just going to repeat what I say?",
	"So immature",
	"Stop copying me!",
}

// repeat concurrently prints out the given message n times
func repeat(n int, message string) {
	var wg sync.WaitGroup

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(m string) {
			log.Println(m)
			wg.Done()
		}(message)
	}
	wg.Wait()
}

func main() {
	factor := flag.Int64("factor", 0, "The fan-out factor to repeat by")
	flag.Parse()
	for _, m := range messages {
		log.Println(m)
		repeat(int(*factor), m)
	}
}
