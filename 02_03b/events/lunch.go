package events

import (
	"encoding/json"
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

/* base case:
- 1 table & 1 line per table
- 1 course
- 1 server
- 1 consumer
*/

/*2 wait groups:
- One for wait staff to deliver orders
- One for clients to take meals
*/

type LunchPlan struct {
	Diners          int      `json:"num_diners"`
	DinerStaffRatio int      `json:"diner_staff_ratio"`
	TableCount      int      `json:"num_tables"`
	Courses         []string `json:"food_courses"`
}

func NewLunchPlanJSON(j json.RawMessage) *LunchPlan {
	var p LunchPlan
	json.Unmarshal(j, &p)
	return &p
}

func (p LunchPlan) Review() {
	log.Printf("p=%#v", p)
}

type Lunch struct {
	diners  int
	staff   int
	courses []string
	tables  []table
}

func (l *Lunch) Plan(p LunchPlan) {
	l.diners = p.Diners
	l.staff = (l.diners / p.DinerStaffRatio) + 1 //from the interwebs
	l.courses = p.Courses

	for ii := range p.TableCount {
		tbl := table{
			num: ii,
		}
		l.tables = append(l.tables, tbl)
		l.tables[ii].setup(ii+1, l.courses)
	}
}

func (l *Lunch) Begin() {
	// Start lunch
	log.Printf("Welcome to the conference lunch! Serving %d attendees.\n",
		l.diners)

	log.Printf("l.diners: %d, l.staff: %d, l.courses: %d", l.diners, l.staff, len(l.courses))

	// Perform server activities
	// Fill 'orders' queue
	orders := l.createOrderQueue()

	var sg sync.WaitGroup
	for range l.staff {
		for o := range orders {
			sg.Go(func() {
				l.tables[o].serveLunch()
			})
		}
	}

	// Perform consumer activities
	// Fill 'meals' queue
	meals := l.createOrderQueue()

	var cg sync.WaitGroup
	for range min(l.diners, len(l.courses)) {
		for m := range meals {
			cg.Go(func() {
				l.tables[m].takeLunch()
			})
		}
	}

	// Both servers and clients must finish their work before the function exits
	sg.Wait()
	cg.Wait()
}

func (l *Lunch) End() {
	//Lunch is over!
	log.Println("Lunch is over!")
}

func (l *Lunch) createOrderQueue() <-chan int {
	jobs := make(chan int) // stream of table numbers to serve/take meals from
	// Fill 'jobs' queue
	go func() {
		for m := range l.diners {
			jobs <- m % len(l.tables)
		}
		close(jobs)
	}()
	return jobs
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

	for _, fc := range t.courses {
		//Use buffer size of 1 to ensure that a server can visit all stations w/o blocking a diner
		t.stations[fc] = &station{
			ready: make(chan struct{}, 1),
		}
	}
}

// takeLunch is the consumer function for the lunch simulation
// Change the signature of this function as required
func (t *table) takeLunch() {
	//A consumer has to visit all stations in order to finish "taking" lunch
	for _, c := range t.courses {
		t.stations[c].take()
		log.Printf("Table: %d, Course: %s, Taken #: %d\n", t.num, c, t.stations[c].taken)
	}
}

// serveLunch is the producer function for the lunch simulation.
// Change the signature of this function as required
func (t *table) serveLunch() {
	//Let a single server deliver an entire lunch to simplify the problem
	log.Printf("Serving lunch at table %d...", t.num)

	for c, s := range t.stations {
		s.serve()
		log.Printf("Table: %d, Course: %s, Served #: %d\n", t.num, c, t.stations[c].served)
	}
}

// Type representing a station where a course is served.
type station struct {
	served int           //Current served
	taken  int           //Current taken
	ready  chan struct{} //Course is available for the next consumer to take
}

func (s *station) serve() {
	// log.Printf("Waiting to serve %d...\n", s.served)
	randomSleep()
	s.ready <- struct{}{}
	s.served++
}

func (s *station) take() {
	// log.Printf("Waiting to take %d...\n", s.taken)
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
