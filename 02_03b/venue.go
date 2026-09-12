package main

import "log"

/*
Type breakdown:
- venue: top level object, containing 1..n tables and 1..m staff, based on the number of consumers
- table: contains all stations required to serve an entire meal
- station: object that produces a course, and signals when it's available for consumption.  Tracks how many meals have been served/consumed
- course: the resource to produce & consume
*/
type venue struct {
	ntables int
	tables  []table
	courses []string
}

func (v *venue) create() {
	for ii := range v.ntables {
		tbl := table{
			num: ii,
		}
		v.tables = append(v.tables, tbl)
		v.tables[ii].setup(ii+1, v.courses)
	}
}

// Type representing a collection of stations representing an entire meal
type table struct {
	num      int                 //Table number
	stations map[string]*station //stations associated with this table
}

// Setup a table with food stations
func (t *table) setup(n int, fcs []string) {
	t.num = n
	t.stations = make(map[string]*station)

	nCourses := 0
	for _, fc := range fcs {
		t.stations[fc] = &station{
			ready: make(chan struct{}, nCourses),
			cap:   consumerCount / nQueues,
		}
	}
}

// Type representing a station where a course is served.
type station struct {
	cap    int           //Max number of servings
	served int           //Current served
	taken  int           //Current taken
	ready  chan struct{} //Course is available for the next consumer to take
}

// Really shouldn't even try to serve after a station is at its capacity...do I need to worry about closing the channel here?
func (s *station) serve() {
	if s.served != s.cap {
		log.Println("Waiting to serve...")
		//randomSleep()
		s.served++
		s.ready <- struct{}{}
	}
}

func (s *station) take() {
	if s.taken != s.cap {
		log.Println("Waiting to take...")
		<-s.ready
		//randomSleep()
		s.taken++
	}
}
