package main

import (
	"encoding/json"
	"io"
	"log"
	"lug/02_03b/events"
	_ "lug/02_03b/events"
	"strings"
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

// Want to be able to create a conference with its events from a JSON file? I don't think I want to support multiple conferences (although I guess I could...)

type Conference []Event

//	func (c *Conference) register(e Event) {
//		c.events = append(c.events, e)
//	}
func NewConference(es string) *Conference {
	dec := json.NewDecoder(strings.NewReader(es))
	p := make(map[string]json.RawMessage)
	for {
		if err := dec.Decode(&p); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
	}

	c := Conference{}
	for k, v := range p {
		switch k {
		case "lunch":
			lp := events.NewLunchPlanJSON(v)
			lp.Review()
			e := events.Lunch{}
			e.Plan(*lp)
			c = append(c, &e)

		default:
			log.Fatalf("Unimplemented event type: %q", k)
		}
	}
	return &c
}

func (c Conference) open() {
	log.Printf("Starting conference with %d events!", len(c))

	for _, e := range c {
		e.Begin()
		e.End()
	}
}

func (c Conference) close() {
	log.Println("Conference is over, thanks for coming!")
}

type Event interface {
	Begin()
	End()
}
