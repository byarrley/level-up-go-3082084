package main

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"
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
const consumerCount = 100
const nQueues = 1    //Serving a buffet lunch to 300 people would take _forever_ with a single queue...but I gotta start somewhere
const nQueueSize = 1 //At most, a queue can have len(foodCourses) diners actively taking food (everyone else is just waiting to start)

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

// takeLunch is the consumer function for the lunch simulation
// Change the signature of this function as required
func takeLunch(t *table, meal int) {
	//A consumer has to visit all stations in order to finish "taking" lunch
	for _, c := range foodCourses {
		t.stations[c].take()
		log.Printf("Meal: %d, Table: %d, Course: %s, Taken #: %d\n", meal, t.num, c, t.stations[c].taken)
	}
}

// serveLunch is the producer function for the lunch simulation.
// Change the signature of this function as required
func serveLunch(t *table, order int) {
	//Let a single server deliver an entire lunch to simplify the problem
	log.Printf("Serving lunch at table %d...", t.num)

	for c, s := range t.stations {
		//Because the courses can be served in any order, if 's.serve()' blocks, it's possible that there will be no consumer
		//  waiting for that course, so the program will deadlock
		go s.serve()
		log.Printf("Order: %d, Table: %d, Course: %s, Served #: %d\n", order, t.num, c, t.stations[c].served)
	}
}

func main() {
	log.Printf("Welcome to the conference lunch! Serving %d attendees.\n",
		consumerCount)

	/*Start with base case:
	- 1 venue (assumed)
	- 1 table
	- 1 line per table
	- 1 course
	- 1 server
	- 1 consumer
	*/

	/*2 wait groups:
	- One for wait staff
	- One for clients
	*/

	// Prepare the venue
	v := venue{ntables: 1,
		courses: foodCourses}
	v.create()

	tbl := &v.tables[0]
	log.Printf("consumerCount: %d, v.courses: %d", consumerCount, len(v.courses))
	fmt.Printf("tbl=%v\n", v.tables[0])

	// Perform server activities
	var sg sync.WaitGroup
	orders := make(chan int) // to be served
	// Fill 'orders' queue
	go func() {
		for m := range consumerCount {
			orders <- m
		}
		close(orders)
	}()

	for range serverCount {
		for o := range orders {
			sg.Go(func() {
				serveLunch(tbl, o)
			})
		}
	}

	// Perform consumer activities
	var cg sync.WaitGroup

	// Fill 'meals' queue
	meals := make(chan int) // to be taken
	go func() {
		for m := range consumerCount {
			meals <- m
		}
		close(meals)
	}()
	for range min(consumerCount, len(v.courses)) {
		for m := range meals {
			cg.Go(func() {
				takeLunch(tbl, m)
			})
		}
	}

	// Both servers and clients must finish their work before the program exits
	sg.Wait()
	cg.Wait()
}

// Sleep to represent activities; maybe create separate functions for serve/consume actions?
func randomSleep() {
	const maxSeconds = 1
	r := rand.Intn(maxSeconds)
	time.Sleep(time.Duration(r)*time.Second + 500*time.Millisecond)
}
