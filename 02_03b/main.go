package main

import (
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
const consumerCount = 2

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
	/*Start with base case:
	- 1 conference (assumed)
	- 1 event (lunch)
	*/

	// Prepare the events
	l := events.Lunch{Ntables: 2,
		Courses: foodCourses}

	//Create the conference and register each event
	c := Conference{attendees: consumerCount}
	c.register(&l)

	//Start the conference
	c.open()

	//End the conference
	c.close()
}
