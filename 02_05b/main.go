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
*/

// setup constants
const baristaCount = 1
const customerCount = 1
const maxOrderCount = 1

// the total amount of drinks that the bartenders have made
type coffeeShop struct {
	orderCount int

	orderCoffee  chan struct{}
	finishCoffee chan struct{}

	//Fan-out(?) orders from a queue
	nextCustomer chan struct{}
	closeShop    chan struct{}

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
		case <-p.orderCoffee:
			p.registerOrder()
			log.Printf("%s makes a coffee.\n", name)
			p.finishCoffee <- struct{}{}
		case <-p.closeShop:
			log.Printf("%s leaves\n", name)
			return
		}
	}
}

// customer is the resource consumer of the coffee shop
func (p *coffeeShop) customer(name string) {
	for {
		select {
		case <-p.nextCustomer:
			p.orderCoffee <- struct{}{}
			log.Printf("%s orders a coffee!\n", name)
			<-p.finishCoffee
			log.Printf("%s enjoys a coffee!\n", name)
		case <-p.closeShop:
			log.Printf("%s leaves\n", name)
			return
		}
	}
}

func main() {
	log.Println("Welcome to the Level Up Go coffee shop!")
	orderCoffee := make(chan struct{}, baristaCount)
	finishCoffee := make(chan struct{}, baristaCount)
	nextCustomer := make(chan struct{})
	closeShop := make(chan struct{})

	p := coffeeShop{
		orderCoffee:  orderCoffee,
		finishCoffee: finishCoffee,
		nextCustomer: nextCustomer,
		closeShop:    closeShop,
	}

	//The shop won't take more than maxOrderCount orders, and due to the simplifications, we can treat this like a work queue and close the channel when all jobs have been submitted
	go func() {
		for range maxOrderCount {
			nextCustomer <- struct{}{}
		}
		close(nextCustomer)
	}()

	for i := 0; i < baristaCount; i++ {
		go p.barista(fmt.Sprint("Barista-", i))
	}
	for i := 0; i < customerCount; i++ {
		go p.customer(fmt.Sprint("Customer-", i))
	}
	<-closeShop

	//This line introduces a data race, because p.orderCount is still being updated after the shop was "closed"!
	log.Printf("Total coffees served: %d", p.orderCount)
	log.Println("The Level Up Go coffee shop has closed! Bye!")
}
