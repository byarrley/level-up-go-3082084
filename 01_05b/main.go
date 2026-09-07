package main

import (
	"encoding/json"
	"flag"
	"log"
	"math"
	"os"
	"sort"
)

const path = "items.json"

// SaleItem represents the item part of the big sale.
type SaleItem struct {
	Name           string  `json:"name"`
	OriginalPrice  float64 `json:"originalPrice"`
	ReducedPrice   float64 `json:"reducedPrice"`
	SalePercentage float64
}

// matchSales adds the sales procentage of the item
// and sorts the array accordingly.
func matchSales(budget float64, items []SaleItem) []SaleItem {
	// Compute % off for each item
	// Sort 'items' in descending order by % off
	// Return the slice of items with the largest discounts that are within budget

	// Solution video: was just supposed to loop over the items, display whatever was in the budget, sorted by SalesPercentage in decreasing order...not figure out what to buy ><

	// Apparently 'range' returns a copy of each slice item, so use the index instead
	for idx := range items {
		items[idx].SalePercentage = math.Round(100 * (items[idx].OriginalPrice - items[idx].ReducedPrice) / items[idx].OriginalPrice)
	}

	// Using ">" for the "less" function so we don't have to reverse the slice
	sort.Slice(items, func(i, j int) bool {
		return items[i].SalePercentage > items[j].SalePercentage
	})

	var list []SaleItem

	for _, i := range items {
		if budget > i.ReducedPrice {
			list = append(list, i)
			budget -= i.ReducedPrice
		}
	}
	return list
}

func main() {
	budget := flag.Float64("budget", 0.0,
		"The max budget you want to shop with.")
	flag.Parse()
	items := importData()
	matchedItems := matchSales(*budget, items)
	printItems(matchedItems)
}

// printItems prints the items and their sales.
func printItems(items []SaleItem) {
	log.Println("The BIG sale has started with our amazing offers!")
	if len(items) == 0 {
		log.Println("No items found.:( Try increasing your budget.")
	}
	for i, r := range items {
		log.Printf("[%d]:%s is %.2f OFF! Get it now for JUST %.2f!\n",
			i, r.Name, r.SalePercentage, r.ReducedPrice)
	}
}

// importData reads the raffle entries from file and
// creates the entries slice.
func importData() []SaleItem {
	file, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	var data []SaleItem
	err = json.Unmarshal(file, &data)
	if err != nil {
		log.Fatal(err)
	}

	return data
}
