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
func (a *auctioneer) runAuction() {

	for _, item := range items {
		bc := make(chan *bid)
		log.Printf("Opening bids for %s!\n", item)

		var bids []bid = make([]bid, 0)
		//Each bidder places their bid on the channel (bc)
		go func() {
			for _, b := range a.bidders {
				b.placeBid(bc)
			}
			close(bc)
		}()

		//For each bid, print it and add it to the slice, then sort by amount
		for bid := range bc {
			// log.Printf("bid=%v\n", bid)
			bids = append(bids, *bid)
		}
		slices.SortStableFunc(bids, func(a, b bid) int {
			return cmp.Compare(a.amount, b.amount)
		})

		//wbid = winning bid
		wbid := bids[len(bids)-1]
		a.bidders[wbid.bidderID].payBid(wbid.amount)
		a.bidders[wbid.bidderID].won++
		log.Printf("%s won %s for $%d! Wallet remaining: $%d, Won: %d\n", wbid.bidderID, item, wbid.amount, a.bidders[wbid.bidderID].wallet, a.bidders[wbid.bidderID].won)
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
func (b *bidder) placeBid(bc chan<- *bid) {
	mybid := bid{
		bidderID: string(b.id),
		amount:   getRandomAmount(b.wallet),
	}
	bc <- &mybid
}

// payBid subtracts the bid amount from the wallet of the auction winner
func (b *bidder) payBid(amount int) {
	b.wallet -= amount
}

func main() {
	rand.Seed(time.Now().UnixNano())
	log.Println("Welcome to the LinkedIn Learning auction.")
	bidders := make(map[string]*bidder, bidderCount)
	for i := 0; i < bidderCount; i++ {
		id := fmt.Sprint("Bidder ", i)
		b := bidder{
			id:     id,
			wallet: walletAmount,
		}
		bidders[id] = &b
	}
	a := auctioneer{
		bidders: bidders,
	}
	a.runAuction()
	log.Println("The LinkedIn Learning auction has finished!")
}

// getRandomAmount generates a random integer amount up to max
func getRandomAmount(max int) int {
	return rand.Intn(int(max))
}
