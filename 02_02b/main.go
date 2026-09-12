package main

import (
	"log"
	"math/rand"
	"time"
)

// The task: Given a list of actions with random durations, implement a function that simulates the concurrent execution
// of the ordered list of actions
//
// Plan:
//   - Assume that each list of actions has to be performed in order
//   - Assume that each list of actions can be performed concurrently with other lists...there are some missing actions that are dependent on the other list
//     but I'm assuming that the lists are simplified to keep the challenge manageable
//   - Try using channels to signal when a list is complete
//
// Solution:
//   - Instructor used sync.WaitGroup
const maxSeconds = 3

type Dog struct {
	name string
}

func (d Dog) fetchLeash() {
	log.Printf("%s goes to fetch leash.\n", d.name)
	randomSleep()
	log.Printf("%s has fetched leash. Woof woof!\n", d.name)
}

func (d Dog) findTreats() {
	log.Printf("%s goes to fetch treats.\n", d.name)
	randomSleep()
	log.Printf("%s has fetched the treats. Woof woof!\n", d.name)
}

func (d Dog) runOutside() {
	log.Printf("%s starts running outside.\n", d.name)
	randomSleep()
	log.Printf("%s is having fun outside. Woof woof!\n", d.name)
}

type Owner struct {
	name string
}

func (o Owner) putShoesOn() {
	log.Printf("%s starts putting shoes on.\n", o.name)
	randomSleep()
	log.Printf("%s finishes putting shoes on.\n", o.name)
}

func (o Owner) findKeys() {
	log.Printf("%s starts looking for keys.\n", o.name)
	randomSleep()
	log.Printf("%s has found keys.\n", o.name)
}

func (o Owner) lockDoor() {
	log.Printf("%s starts locking the door.\n", o.name)
	randomSleep()
	log.Printf("%s has locked the door.\n", o.name)
}

func randomSleep() {
	r := rand.Intn(maxSeconds)
	time.Sleep(time.Duration(r)*time.Second + 500*time.Millisecond)
}

func main() {
	owner := Owner{name: "Jimmy"}
	dog := Dog{name: "Lucky"}
	ownerActions := []func(){
		owner.putShoesOn,
		owner.findKeys,
		owner.lockDoor,
	}
	dogActions := []func(){
		dog.fetchLeash,
		dog.findTreats,
		dog.runOutside,
	}

	jobs := make(chan []func(), 2)
	jobs <- ownerActions
	jobs <- dogActions
	close(jobs)

	executeWalk(jobs)
}

func executeWalk(actionList <-chan []func()) {
	//This may not be guaranteed to work if actionList is an unbuffered channel
	nList := len(actionList)
	done := make(chan struct{})

	for a := range actionList {
		go doActions(a, done)
	}

	nDone := 0
	for range done {
		nDone++

		if nDone == nList {
			//Normally, the receiver wouldn't close a channel, but since we know exactly how many responses to expect,
			//  and the channel is local to the function, it seems safe to do here?
			close(done)
		}
	}
}

// Process a list of actions
func doActions(actions []func(), done chan<- struct{}) {
	for _, a := range actions {
		a()
	}
	done <- struct{}{}
}
