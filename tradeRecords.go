package main

/*
	Hello Messari team members, thank you for considering me for this position, and taking the
	time to review my work.

	A quick disclaimer: Before this project, my only experience with Go was the hello
	world I ran about five minuetes before creating this file, regardless, the language has really
	grown on me.

	That being said, this was definitely one of the better take home projects that I've recieved.

	I look forward to hearing back from you!

*/

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"messariTradeRecords/funcs" // additional functions and types located in funcs/helperFunctions.go
	"os"
	"time"

	"github.com/montanaflynn/stats" // stats library for calculations
)

func main() {

	/*
		The bufio scanner provided the easiest way to read both the text and Json outputs from the binary.
		However, during my inital I found the scanner was unable to buffer the required space to read all the input,
		and so is manually buffered at 100 as input size ranges from 94-97 bytes.
	*/

	scanner := bufio.NewScanner(os.Stdin)
	buffer := make([]byte, 100)
	scanner.Buffer(buffer, 100)

	// Initialize RecMap
	r := funcs.RecMap{}

	for scanner.Scan() {
		// Trade type
		var t funcs.Trade
		// Initialize Pesist
		p := funcs.Persist{}

		// switch statement on scanner Text to read BEGIN/END.
		switch scanner.Text() {

		case "BEGIN":
			fmt.Println("Consuming Data:")
		case "END":
			fmt.Println("Data Consumed")
			break
		default:
			// Unmarshal bytes and log Fatal error.
			x := json.Unmarshal(scanner.Bytes(), &t)

			if x != nil {
				log.Fatal(x)
				continue
			}
			// Aggregate trade data.
			i := funcs.Agg(t, p, r)

			r[t.Market] = i
		}
	}
	start := time.Now()
	for key, element := range r {

		// Iterate through map of market data on read finish, populating the Record struct for final Json marshal.

		fin := funcs.Record{}
		trades := len(element.Prices)
		meanP, _ := stats.Mean(element.Prices)
		meanV, _ := stats.Mean(element.Volume)
		totV, _ := stats.Sum(element.Volume)
		pBuy, _ := stats.Sum(element.Percentage_buy)

		fin.Market = key // Market
		fin.TotalVolume = funcs.RoundFloat(totV, 2)
		fin.MeanPrice = funcs.RoundFloat(meanP, 2)
		fin.MeanVolume = funcs.RoundFloat(meanV, 2)
		fin.VWAP = funcs.RoundFloat(element.VWAP/totV, 2)
		fin.Percentage_buy = funcs.RoundFloat((pBuy / float64(trades)), 2)

		aggData, _ := json.Marshal(fin)

		fmt.Println(string(aggData))

	}
	end := time.Now()
	fmt.Println("Took", (end.Sub(start)), "to aggregate the data.")

}
