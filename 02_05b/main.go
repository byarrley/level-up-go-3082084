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
*/

// setup constants
const baristaCount = 3
const customerCount = 20
const maxOrderCount = 40

// the total amount of drinks that the bartenders have made
type coffeeShop struct {
	orderCount int

	orderCoffee  chan struct{}
	finishCoffee chan struct{}
	closeShop    chan struct{}

	mux sync.Mutex
}

// registerOrder ensures that the order made by the baristas is counted
func (p *coffeeShop) registerOrder() {
	p.mux.Lock()
	defer p.mux.Unlock()

	p.orderCount++
	if p.orderCount == maxOrderCount {
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
			log.Printf("%s leaves", name)
			return
		}
	}
}

// customer is the resource consumer of the coffee shop
func (p *coffeeShop) customer(name string) {
	for {
		select {
		case p.orderCoffee <- struct{}{}:
			log.Printf("%s orders a coffee!", name)
			<-p.finishCoffee
			log.Printf("%s enjoys a coffee!\n", name)
		case <-p.closeShop:
			log.Printf("%s leaves", name)
			return
		}
	}
}

func main() {
	log.Println("Welcome to the Level Up Go coffee shop!")
	orderCoffee := make(chan struct{}, baristaCount)
	finishCoffee := make(chan struct{}, baristaCount)
	closeShop := make(chan struct{})
	p := coffeeShop{
		orderCoffee:  orderCoffee,
		finishCoffee: finishCoffee,
		closeShop:    closeShop,
	}
	for i := 0; i < baristaCount; i++ {
		go p.barista(fmt.Sprint("Barista-", i))
	}
	for i := 0; i < customerCount; i++ {
		go p.customer(fmt.Sprint("Customer-", i))
	}
	<-p.closeShop

	log.Printf("Total coffees served: %d", p.orderCount)
	log.Println("The Level Up Go coffee shop has closed! Bye!")
}
