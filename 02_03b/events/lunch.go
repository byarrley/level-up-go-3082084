package events

import (
	"log"
	"math/rand"
	"sync"
	"time"
)

/*
- lunch (event): containing 1..n tables and 1..m staff, based on the number of consumers
- table: contains all stations required to serve an entire meal
- station: object that produces a course, and signals when it's available for consumption.  Tracks how many meals have been served/consumed
- course: the resource to produce & consume
*/

type Lunch struct {
	Ntables int
	staff   int
	diners  int
	tables  []table
	Courses []string
}

func (l *Lunch) Plan(attendees int) {
	l.diners = attendees
	l.staff = (attendees / 50) + 1 //from the interwebs
	for ii := range l.Ntables {
		tbl := table{
			num: ii,
		}
		l.tables = append(l.tables, tbl)
		l.tables[ii].setup(ii+1, l.Courses)
	}
}

func (l *Lunch) Begin() {
	// Start lunch
	log.Printf("Welcome to the conference lunch! Serving %d attendees.\n",
		l.diners)

	log.Printf("consumerCount: %d, l.Courses: %d", l.diners, len(l.Courses))
	// fmt.Printf("tbl=%v\n", l.tables[0])

	// Perform server activities
	var sg sync.WaitGroup
	orders := make(chan int) // to be served
	// Fill 'orders' queue
	go func() {
		for m := range l.diners {
			orders <- m % len(l.tables)
		}
		close(orders)
	}()

	for range l.staff {
		for o := range orders {
			sg.Go(func() {
				l.tables[o].serveLunch(o)
			})
		}
	}

	// Perform consumer activities
	var cg sync.WaitGroup

	// Fill 'meals' queue
	meals := make(chan int) // to be taken
	go func() {
		for m := range l.diners {
			meals <- m % len(l.tables)
		}
		close(meals)
	}()
	for range min(l.diners, int(len(l.Courses))) {
		for m := range meals {
			cg.Go(func() {
				l.tables[m].takeLunch(m)
			})
		}
	}

	// Both servers and clients must finish their work before the program exits
	sg.Wait()
	cg.Wait()
}

func (l *Lunch) End() {
	//Lunch is over!
	log.Println("Lunch is over!")
}

// Type representing a collection of stations representing an entire meal
type table struct {
	num      int                 //Table number
	stations map[string]*station //stations associated with this table
	courses  []string            //Courses that make up the meal served at this table
}

// Setup a table with food stations
func (t *table) setup(n int, fcs []string) {
	t.num = n
	t.stations = make(map[string]*station)
	t.courses = fcs

	nCourses := 0
	for _, fc := range t.courses {
		t.stations[fc] = &station{
			ready: make(chan struct{}, nCourses),
		}
	}
}

// takeLunch is the consumer function for the lunch simulation
// Change the signature of this function as required
func (t *table) takeLunch(meal int) {
	//A consumer has to visit all stations in order to finish "taking" lunch
	for _, c := range t.courses {
		t.stations[c].take()
		log.Printf("Meal: %d, Table: %d, Course: %s, Taken #: %d\n", meal, t.num, c, t.stations[c].taken)
	}
}

// serveLunch is the producer function for the lunch simulation.
// Change the signature of this function as required
func (t *table) serveLunch(order int) {
	//Let a single server deliver an entire lunch to simplify the problem
	log.Printf("Serving lunch at table %d...", t.num)

	for c, s := range t.stations {
		//Because the courses can be served in any order, if 's.serve()' blocks, it's possible that there will be no consumer
		//  waiting for that course, so the program will deadlock
		go s.serve()
		log.Printf("Order: %d, Table: %d, Course: %s, Served #: %d\n", order, t.num, c, t.stations[c].served)
	}
}

// Type representing a station where a course is served.
type station struct {
	served int           //Current served
	taken  int           //Current taken
	ready  chan struct{} //Course is available for the next consumer to take
}

func (s *station) serve() {
	log.Printf("Waiting to serve %d...\n", s.served)
	randomSleep()
	s.ready <- struct{}{}
	s.served++
}

func (s *station) take() {
	log.Printf("Waiting to take %d...\n", s.taken)
	randomSleep()
	<-s.ready
	s.taken++
}

// Sleep to represent activities; maybe create separate functions for serve/consume actions?
func randomSleep() {
	const maxSeconds = 1
	r := rand.Intn(maxSeconds)
	time.Sleep(time.Duration(r)*time.Second + 500*time.Millisecond)
}
