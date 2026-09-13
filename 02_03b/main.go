package main

import (
	"log"
	"lug/02_03b/events"
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
- Each attendee must wait to take a course until the previous attendee has finished taking it (mutex?)
- Each attendee takes exactly one of each course
- Travel and consumption time are excluded from the simulation
*/

// the number of attendees we need to serve lunch to
const consumerCount = 1

// servers
const serverCount = (consumerCount / 50) + 1 //from the interwebs

// foodCourses represents the types of resources to pass to the consumers
// var foodCourses = []string{
// 	"Caprese Salad",
// 	"Spaghetti Carbonara",
// 	"Vanilla Panna Cotta",
// }

var foodCourses = []string{
	"Caprese Salad",
}

func main() {
	log.Printf("Welcome to the conference lunch! Serving %d attendees.\n",
		consumerCount)

	/*Start with base case:
	- 1 conference (assumed)
	- 1 event (lunch)
	- 1 table & 1 line per table
	- 1 course
	- 1 server
	- 1 consumer
	*/

	/*2 wait groups:
	- One for wait staff
	- One for clients
	*/

	// Prepare the event
	l := events.Lunch{Ntables: 1,
		Courses: foodCourses}
	l.Plan(consumerCount)
	l.Begin()
	l.End()

}
