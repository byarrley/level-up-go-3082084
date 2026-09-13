package main

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

type Event interface {
	Plan(attendees int)
	Begin()
	End()
}
