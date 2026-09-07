package main

import (
	"flag"
	"fmt"
	"log"
	"math"
)

// coin contains the name and value of a coin
type coin struct {
	name  string
	value float64
}

// coins is the list of values available for making change.
var coins = []coin{
	{name: "1 pound", value: 1},
	{name: "50 pence", value: 0.50},
	{name: "20 pence", value: 0.20},
	{name: "10 pence", value: 0.10},
	{name: "5 pence", value: 0.05},
	{name: "1 penny", value: 0.01},
}

// calculateChange returns the coins required to calculate the
func calculateChange(amount float64) map[coin]int {
	chg := make(map[coin]int, 0)

	for _, c := range coins {
		//Solution video: can use division to get a count of each coin instead of the 2nd inner loop
		for {
			if c.value > amount {
				break
			}
			chg[c]++
			amount -= c.value
			//Have to round, or something like 1.76 will be off by 1 penny
			amount = math.Round(amount*100) / 100
			fmt.Println(amount)
		}
	}
	return chg
}

// printCoins prints all the coins in the slice to the terminal.
func printCoins(change map[coin]int) {
	if len(change) == 0 {
		log.Println("No change found.")
		return
	}
	log.Println("Change has been calculated.")
	for coin, count := range change {
		log.Printf("%d x %s \n", count, coin.name)
	}
}

func main() {
	amount := flag.Float64("amount", 0.0, "The amount you want to make change for")
	flag.Parse()
	change := calculateChange(*amount)
	printCoins(change)
}
