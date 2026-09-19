package main

import (
	"fmt"
	"log"
	"sync"
)

// The Task: Given a faulty concurrent simulation, implement a fix to ensure that there are no race conditions or crashes
/* Simplifications:
- We will use fulfilled orders as a proxy to passing time
- We will not care about which customer has ordered which coffee
- We will not maintain any coffee types; all orders will be represented by the empty struct - struct{}

Hints:
- Run the program using the '-race' flag to detect bugs.  Make use of stack traces on potential crashes as well.

Solution:
- I got stuck on this one, it turns out that the solution is to use a third channel that signals to the other goroutines that the shop is closed, by simply closing the channel
- However, I'm curious to see if the order count goes over the max before the program actually shuts down

Post-solution notes:
- Yes, the count _can_ go over the max before the program exits (if you run it enough times)
- If it's possible with having the same channel send in the consumer select and receive in the barista select, I don't see it.
	There are few examples of send channels in use (at least that I've been able to find), and they tend to be for time outs.

Command line tests:
- for ii in {1..10}; do echo "*** run $ii ***"; go run -race main.go 2>&1 | tee runs/$ii.log ; sleep 1; echo; done # Run 10x, tee results to file
- grep -nE '(leaves|clocks|shortly|Time|Bye)' # Print lines tracking when customers & baristas leave the store relative to their respective announcements
*/

// setup constants
const baristaCount = 2
const customerCount = 8
const maxOrderCount = 8

// the total amount of drinks that the bartenders have made
type coffeeShop struct {
	orderCount int

	orderCoffee  chan string
	finishCoffee chan string

	nextCustomer chan struct{} //Fan-out(?) orders from a queue

	closeShop chan struct{} //Signal to customers that the shop is closed

	//Channels to signal which baristas/customers have left.  If the status isn't sent with the signal, the logs may become out of sync with the signals
	baristaLeft  chan string
	customerLeft chan string

	mux sync.Mutex
}

// registerOrder ensures that the order made by the baristas is counted
func (p *coffeeShop) registerOrder() {
	p.mux.Lock()
	defer p.mux.Unlock()

	p.orderCount++
	if p.orderCount == maxOrderCount {
		//This is a neat trick from the solution to use a channel signal that a process is finished without sending anything over it
		close(p.closeShop)
	}
}

// barista is the resource producer of the coffee shop
func (p *coffeeShop) barista(name string) {
	for {
		select {
		//p.orderCoffee is never closed, so we don't have to ensure that it's open before receiving from it
		case msg := <-p.orderCoffee: //barista: order received
			p.registerOrder()                                             //barista: register order
			log.Printf("%s", msg)                                         //barista: log order
			p.finishCoffee <- fmt.Sprintf("> %s makes a coffee.\n", name) //barista: serve order
		case <-p.closeShop: //barista: receive signal that it's time to clock out
			p.baristaLeft <- fmt.Sprintf("%s clocks out\n", name) //barista: signal that this barista has clocked out
			return
		}
	}
}

// customer is the resource consumer of the coffee shop
func (p *coffeeShop) customer(name string) {
	for {
		select {
		case _, ok := <-p.nextCustomer: //customer: this customer is next to place an order.  Attempting to set the channel to 'nil' introduces data races
			//Use "ok" to ensure that the channel is open before accepting the next customer's order
			if ok {
				p.orderCoffee <- fmt.Sprintf("%s orders a coffee!\n", name) //customer: order sent
				log.Printf("%s", <-p.finishCoffee)                          //barista: coffee served
				log.Printf("%s enjoys a coffee!\n", name)                   //customer: drink coffee
			}
		case <-p.closeShop: //customer: receive signal that shop is closing
			p.customerLeft <- fmt.Sprintf("%s leaves the shop\n", name) //customer: signal that this customer has left the shop
			return
		}
	}
}

func main() {
	log.Println("Welcome to the Level Up Go coffee shop!")
	p := NewCoffeeShop()

	//The shop won't take more than maxOrderCount orders, and due to the simplifications, we can treat this like a work queue and close the channel when all jobs have been submitted
	go func() {
		for range maxOrderCount {
			p.nextCustomer <- struct{}{} //shop: call next customer ready to order
		}
		close(p.nextCustomer) //shop: signal to all customers that no more orders can be placed.  Setting to 'nil' here introduces a data race
	}()

	for i := 0; i < baristaCount; i++ {
		go p.barista(fmt.Sprint("Barista-", i))
	}
	for i := 0; i < customerCount; i++ {
		go p.customer(fmt.Sprint("Customer-", i))
	}
	<-p.closeShop //shop: wait for signal to all customers that shop is closing before continuing

	//Wait for customers to leave before continuing; it's OK if they leave _before_ the announcement, but they should have finished their actions and left before
	// 	the baristas get the signal to clock out
	log.Println("---The Level Up Go coffee shop is closing shortly...---")
	for range customerCount {
		log.Printf("%s", <-p.customerLeft)
	}

	//Annouce work done and that baristas can leave before continuing
	log.Println("***Time to clock out!***")
	log.Printf("Total coffees served: %d.  Great work team!", p.orderCount)
	for range baristaCount {
		log.Printf("%s", <-p.baristaLeft)
	}

	//Block until all customers/baristas have completed their actions and have left
	log.Println("The Level Up Go coffee shop has closed! Bye!")
}

func NewCoffeeShop() *coffeeShop {
	p := coffeeShop{
		nextCustomer: make(chan struct{}),
		orderCoffee:  make(chan string),
		finishCoffee: make(chan string),
		closeShop:    make(chan struct{}),
		customerLeft: make(chan string),
		baristaLeft:  make(chan string),
	}
	return &p
}
