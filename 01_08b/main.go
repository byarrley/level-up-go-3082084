package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"
)

//The Task: Given a list of nodes, implement a [recursive] function that ensures that all nodes are visited.

const path = "friends.json"

// Friend represents a friend and their connections.
type Friend struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Friends     []string `json:"friends"`
	GossipHeard bool
}

// hearGossip indicates that the friend has heard the gossip.
func (f *Friend) hearGossip() {
	f.GossipHeard = true
	log.Printf("%s has heard the gossip!\n", f.Name)
}

// Friends represents the map of friends and connections
type Friends struct {
	fmap map[string]*Friend
}

// getFriend fetches the friend given an id.
func (f *Friends) getFriend(id string) *Friend {
	return f.fmap[id]
}

// getRandomFriend returns an random friend.
func (f *Friends) getRandomFriend() *Friend {
	rand.Seed(time.Now().Unix())
	id := (rand.Intn(len(f.fmap)-1) + 1) * 100
	return f.getFriend(fmt.Sprint(id))
}

// spreadGossip ensures that all the friends in the map have heard the news
func spreadGossip(root *Friend, friends Friends) {

	//Solution video: the instructor modified this function to take another argument, which is a map of visited nodes.
	//								I was operating under the assumption that the provided functionality was complete, or I would've kept trying to make the embedded boolean work...I'll try that next.
	//Attempt 2: I added a bool to the Friend struct indicating whether or not the friend has heard the gossip.  If not, they hear and then spread it.  If so, no-op.
	//					 This seemed like a more intuitive way to do it, but it required passing *Friend so that modifications to the Friends members would affect the originals (not copies)

	//typical case: current friend hears the gossip and then spreads it
	//base case: everyone has heard it (return without doing anything)
	for _, id := range root.Friends {
		f := friends.getFriend(id)

		if !f.GossipHeard {
			f.hearGossip()
			spreadGossip(f, friends)
		}
	}
	//panic("NOT IMPLEMENTED")
}

func main() {
	friends := importData()
	root := friends.getRandomFriend()
	root.hearGossip()
	spreadGossip(root, friends)
}

// importData reads the input data from file and
// creates the friends map.
func importData() Friends {
	file, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	var data []Friend
	err = json.Unmarshal(file, &data)
	if err != nil {
		log.Fatal(err)
	}

	fm := make(map[string]*Friend, len(data))
	for _, d := range data {
		fm[d.ID] = &d
	}

	return Friends{
		fmap: fm,
	}
}
