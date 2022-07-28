package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"github.com/montanaflynn/stats"
)

type trade struct {
	Id     json.Number `json:"id"`
	Market json.Number `json:"market" type:"integer" required:"true"`
	Price  json.Number `json:"price" type:"float64" required:"true"`
	Volume json.Number `json:"volume"`
	Is_buy bool        `json:"is_buy"`
}

type record struct {
	Market         json.Number `type:"integer" required:"true"`
	TotalVolume    float64
	MeanPrice      float64
	MeanVolume     float64
	VWAP           float64
	Percentage_buy float64
}

type recTest map[json.Number][]float64

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	s := recTest{}
	//test := make([]record, 5)
	for scanner.Scan() {
		var t trade

		switch scanner.Text() {

		case "BEGIN":
			fmt.Println("Market Open:")
		case "END":
			fmt.Println("Market Close")
			break
		default:

			x := json.Unmarshal(scanner.Bytes(), &t)
			//.Decode(&t)
			//r := &rec
			// new func for routune?
			if x != nil {
				fmt.Printf("%+v %v", t, x)
			}

			i, _ := t.Price.Float64()
			j, _ := t.Volume.Float64()

			s[t.Market] = append(s[t.Market], i)
			s[t.Market] = append(s[t.Volume], j)
		}
	}

	for key, element := range s {
		trades := len(element)
		mean, _ := stats.Mean(element)
		fmt.Println("Market:", key, " Mean Price:", mean, "Trades: ", trades)

	}
}
