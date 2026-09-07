package main

import (
	"flag"
	"log"
	"math"
	"time"
)

var expectedFormat = "2006-01-02"

// parseTime validates and parses a given date string.
func parseTime(target string) time.Time {
	t, err := time.Parse(expectedFormat, target)

	if err != nil {
		log.Fatalln(err.Error())
	}

	return t
	//panic("NOT IMPLEMENTED")
}

// calcSleeps returns the number of sleeps until the target.
func calcSleeps(target time.Time) float64 {
	durationUntil := time.Until(target)
	return math.Ceil(durationUntil.Hours() / 24.0)
	//panic("NOT IMPLEMENTED")
}

func main() {
	bday := flag.String("bday", "", "Your next bday in YYYY-MM-DD format")
	flag.Parse()
	target := parseTime(*bday)
	log.Printf("You have %d sleeps until your birthday. Hurray!",
		int(calcSleeps(target)))
}
