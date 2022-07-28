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
	Percentage_buy float64
	VWAP           float64
}

type recTest map[json.Number]persist

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	s := recTest{}
	//test := make([]record, 5)
	for scanner.Scan() {
		var t trade
		m := persist{}

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
				fmt.Println(x)
				continue
			}

			i, _ := t.Price.Float64()
			j, _ := t.Volume.Float64()

			m.Prices = append(m.Prices, i)
			m.Volume = append(m.Volume, j)
			sumvol, _ := stats.Sum(m.Volume)
			m.VWAP = (i * j) / sumvol

			if t.Is_buy == true {
				m.Percentage_buy += 1.0
			}
			s[t.Market] = m
		}
	}
	//fmt.Println(s)

	for key, element := range s {
		fin := record{}
		trades := len(element.Prices)
		meanP, _ := stats.Mean(element.Prices)
		meanV, _ := stats.Mean(element.Volume)
		totV, _ := stats.Sum(element.Volume)

		fin.Market = key
		fin.TotalVolume = totV
		fin.MeanPrice = meanP
		fin.MeanVolume = meanV
		fin.VWAP = element.VWAP
		fin.Percentage_buy = (element.Percentage_buy / float64(trades))

		aggData, _ := json.Marshal(fin)
		fmt.Println(string(aggData))

	}

}
