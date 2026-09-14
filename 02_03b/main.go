package main

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

func main() {
	/*Start with base case:
	- 1 conference (assumed)
	- 1 event (lunch)
	*/

	const eventStream = `
		{
			"lunch": {
				"num_diners": 2,
				"staff_ratio": 50, 
				"num_tables": 2,
				"food_courses": ["Caprese Salad", "Spaghetti Carbonara"]
			}
		}
	`
	c := NewConference(eventStream)

	// //Start the conference
	c.open()

	// //End the conference
	c.close()
}
