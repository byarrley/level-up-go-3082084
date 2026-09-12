package main

import (
	"log"
)

// The Task: Given a defined list of resources, implement a function that simulates the concurrent allocation of resources to consumer goroutines.
// Hint: A signal channel, `var signal chan struct{}`, is a channel whose purpose is to synchronize goroutines

/* Assumptions:
serveLunch
- Each course takes a random amount of time to be served
- Courses are served concurrently
- Assume 6:1 attendee to service staff ratio
- There are enough of each course for all consumers
- Each course has a buffer associated with it
takeLunch
- Each course is received in slice order, requires a random amount of time to be taken
- Buffet-style lunch (each consumer takes all 3 courses in order before they exit the line)
- Line size should be limited so as not to exhaust system resources
- Assume a single line
- Each attendee must wait to take a course until the previous attendee has finished taking it
- Each attendee takes exactly one of each course
*/

// the number of attendees we need to serve lunch to
const consumerCount = 300

// foodCourses represents the types of resources to pass to the consumers
var foodCourses = []string{
	"Caprese Salad",
	"Spaghetti Carbonara",
	"Vanilla Panna Cotta",
}

// takeLunch is the consumer function for the lunch simulation
// Change the signature of this function as required
func takeLunch(name string) {
	panic("NOT IMPLEMENTED YET")
}

// serveLunch is the producer function for the lunch simulation.
// Change the signature of this function as required
func serveLunch(course string) {
	panic("NOT IMPLEMENTED YET")
}

func main() {
	log.Printf("Welcome to the conference lunch! Serving %d attendees.\n",
		consumerCount)
}
