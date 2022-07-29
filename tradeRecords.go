package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math"
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

type recTest map[json.Number]persist

func roundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	buffer := make([]byte, 100)
	scanner.Buffer(buffer, 100)

	s := make(recTest)

	for scanner.Scan() {

		var t trade
		//var b *persist
		b := persist{}

		switch scanner.Text() {

		case "BEGIN":
			fmt.Println("Consuming Data:")
		case "END":
			fmt.Println("Market Close")
			break
		default:

			x := json.Unmarshal(scanner.Bytes(), &t)
			//.Decode(&t)
			//r := &rec
			// new func for routune?
			if x != nil {
				fmt.Println(x)
				continue
			}
			b = s[t.Market]

			p, _ := t.Price.Float64()
			v, _ := t.Volume.Float64()

			b.Prices = append(b.Prices, p)
			b.Volume = append(b.Volume, v)

			max, _ := stats.Max(b.Prices)
			min, _ := stats.Min(b.Prices)
			b.VWAP += ((max + min + p) / 3) * v

			if t.Is_buy == true {
				b.Percentage_buy = append(b.Percentage_buy, float64(1))
			} else {
				b.Percentage_buy = append(b.Percentage_buy, float64(0))
			}
			//fmt.Println(m)
			s[t.Market] = b

		}
	}

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

}
