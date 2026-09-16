package main

import (
	"cmp"
	"fmt"
	"log"
	"math/rand"
	"slices"
	"time"
)

// The Task: Given a list of resources, implement a function that simulates the concurrent management and allocation of bids for those resources
/* Simplifications:
- All bidders start with equal amounts of money
- All bids will be integer amounts
- All bidders will bid for all items
- If a bidder is out of money, they will place a 0 bid
- No minimum prices for our items
- Items are sold one by one, sequentially

Hints: Make use of the mechanisms we are already familiar with.  You've got this!

Plan
- For each item
	* Gather bids from the participants
		- Bids can range from 0 up to bidder.wallet
	* Identify the bidder with the highest bid, deduct the amount from their wallet

Notes
- Not sure why the bidders are placing bids as soon as they're created (without an item being announced), so I will move that call

Solution video:
- The bidders were placing bids because the instructor was expecting to use a signal channel and a bid channel and pass those between the bidders and the auctioneer
- Auctioneer signals to all bidders that the auction is open for the current item, and receives all "bid" objects via the other channel
- Each bidder waits to receive the auction open signal, then sends its bid on the bid channel
- It didn't occur to me that I shouldn't change the main code, but it was possible to solve this challenge just by modifying the function signatures (and creating helper functions)
*/

// the amount of bidders we have at our auction
const bidderCount = 10

// initial wallet value for all bidders
const walletAmount = 250

// items is the map of auction items
var items = []string{
	"The \"Best Gopher\" trophy",
	"The \"Learn Go with Adelina\" experience",
	"Two tickets to a Go conference",
	"Signed copy of \"Beautiful Go code\"",
	"Vintage Gopher plushie",
}

// bid is a type that pairs the bidder id and the amount they want to bid
type bid struct {
	bidderID string
	amount   int
}

// auctioneer receives bids and announces winners
type auctioneer struct {
	bidders map[string]*bidder
}

// runAuction and manages the auction for all the items to be sold
// Change the signature of this function as required
func (a *auctioneer) runAuction(openc chan<- struct{}, bidc <-chan bid) {

	for _, item := range items {
		log.Printf("Opening bids for %s!\n", item)

		//Grab each bid, print it and add it to the slice, then sort by amount
		var bids []bid
		for range bidderCount {
			//Signal to bidder that the auction is open for this item
			openc <- struct{}{}
			//Take bid
			bid := <-bidc
			log.Printf("%s: $%d\n", bid.bidderID, bid.amount)
			bids = append(bids, bid)
		}

		slices.SortStableFunc(bids, func(a, b bid) int {
			return cmp.Compare(a.amount, b.amount)
		})

		//wbid = winning bid
		wbid := bids[len(bids)-1]
		a.bidders[wbid.bidderID].payBid(wbid.amount)
		a.bidders[wbid.bidderID].won++
		log.Printf("%s won %s for $%d! Wallet remaining: $%d, Won: %d\n\n", wbid.bidderID, item, wbid.amount, a.bidders[wbid.bidderID].wallet, a.bidders[wbid.bidderID].won)
	}
}

// bidder is a type that holds the bidder id and wallet
type bidder struct {
	id     string
	wallet int
	won    int
}

// placeBid generates a random amount and places it on the bids channels
// Change the signature of this function as required
func (b *bidder) placeBid(openc <-chan struct{}, bidc chan<- bid) {
	for range len(items) {
		<-openc
		mybid := bid{
			bidderID: string(b.id),
			amount:   getRandomAmount(b.wallet),
		}
		bidc <- mybid
	}
}

// payBid subtracts the bid amount from the wallet of the auction winner
func (b *bidder) payBid(amount int) {
	b.wallet -= amount
}

func main() {
	rand.Seed(time.Now().UnixNano())
	log.Println("Welcome to the LinkedIn Learning auction.")

	//openc: signal channel to tell bidders that auction is open for the current item
	//bidc: buffered channel for bidders to send bids to the auctioneer
	openc := make(chan struct{})
	bidc := make(chan bid, bidderCount)

	bidders := make(map[string]*bidder, bidderCount)
	for i := 0; i < bidderCount; i++ {
		id := fmt.Sprint("Bidder ", i)
		b := bidder{
			id:     id,
			wallet: walletAmount,
		}
		bidders[id] = &b
		go b.placeBid(openc, bidc)
	}
	a := auctioneer{
		bidders: bidders,
	}
	a.runAuction(openc, bidc)
	log.Println("The LinkedIn Learning auction has finished!")
}

// getRandomAmount generates a random integer amount up to max
func getRandomAmount(max int) int {
	return rand.Intn(int(max))
}
