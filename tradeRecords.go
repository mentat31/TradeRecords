package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"time"

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
	Market         json.Number `json:"market"`
	TotalVolume    float64     `json:"Total Volume"`
	MeanPrice      float64     `json:"Mean Price"`
	MeanVolume     float64     `json:"Mean Volume"`
	VWAP           float64     `json:"VWAP"`
	Percentage_buy float64     `json:"Percentage Buy"`
}

type persist struct {
	Market         json.Number
	Volume         []float64
	Prices         []float64
	Percentage_buy []float64
	VWAP           float64
}

type recMap map[json.Number]persist

func roundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}

func agg(current trade, record persist, marketData recMap) persist {

	r := marketData[current.Market]

	p, _ := current.Price.Float64()
	v, _ := current.Volume.Float64()

	r.Prices = append(r.Prices, p)
	r.Volume = append(r.Volume, v)

	max, _ := stats.Max(r.Prices)
	min, _ := stats.Min(r.Prices)
	r.VWAP += ((max + min + p) / 3) * v

	if current.Is_buy == true {
		r.Percentage_buy = append(r.Percentage_buy, float64(1))
	} else {
		r.Percentage_buy = append(r.Percentage_buy, float64(0))
	}
	marketData[current.Market] = r
	return marketData[current.Market]
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	buffer := make([]byte, 100)
	scanner.Buffer(buffer, 100)

	s := recMap{}

	for scanner.Scan() {

		var t trade

		b := persist{}

		switch scanner.Text() {

		case "BEGIN":
			fmt.Println("Consuming Data:")
		case "END":
			fmt.Println("Data Consumed")
			break
		default:
			x := json.Unmarshal(scanner.Bytes(), &t)

			if x != nil {
				fmt.Println(x)
				continue
			}
			i := agg(t, b, s)

			s[t.Market] = i
		}
	}
	start := time.Now()
	for key, element := range s {

		fin := record{}
		trades := len(element.Prices)
		meanP, _ := stats.Mean(element.Prices)
		meanV, _ := stats.Mean(element.Volume)
		totV, _ := stats.Sum(element.Volume)
		pBuy, _ := stats.Sum(element.Percentage_buy)

		fin.Market = key
		fin.TotalVolume = roundFloat(totV, 2)
		fin.MeanPrice = roundFloat(meanP, 2)
		fin.MeanVolume = roundFloat(meanV, 2)
		fin.VWAP = roundFloat(element.VWAP/totV, 2)
		fin.Percentage_buy = roundFloat((pBuy / float64(trades)), 2)

		aggData, _ := json.Marshal(fin)

		fmt.Println(string(aggData))

	}
	end := time.Now()
	fmt.Println("Took", (end.Sub(start)), "to aggregate the data.")

}
