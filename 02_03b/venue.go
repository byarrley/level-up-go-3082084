package main

import (
	"log"
	_ "lug/02_03b/events"
)

/*
Type breakdown:
- conference: top level object, containing an event schedule, and 1..x registered events
  - plan(): set up the conference with number of attendees, event schedule
  - register(): Add an event to the schedule
  - open(): start the conference with opening remarks
  - close(): stop the conference with closing remarks

- event: interface with the following methods
  - plan(): configure an event, including its staff pool based on the number of attendees, etc
  - begin(): start an event with a message
  - end(): terminate an event with a message, maybe support a timeout?
*/

type Conference struct {
	events    []Event
	attendees int
}

func (c *Conference) register(e Event) {
	c.events = append(c.events, e)
}

func (c *Conference) open() {
	log.Printf("Starting conference with %d events!", len(c.events))

	for _, e := range c.events {
		e.Plan(c.attendees)
		e.Begin()
		e.End()
	}
}

func (c *Conference) close() {
	log.Println("Conference is over, thanks for coming!")
}

type Event interface {
	Plan(attendees int)
	Begin()
	End()
}
