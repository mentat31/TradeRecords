package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type trade struct {
	Id     json.Number `json:"id"`
	Market json.Number `json:"market"`
	Price  json.Number `json:"price" type:"float" required:"true"`
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

type recTest map[json.Number][]json.Number

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
			break
		default:

			x := json.NewDecoder(strings.NewReader(scanner.Text())).Decode(&t)
			//r := &rec
			// new func for routune?
			if x != nil {
				fmt.Println(x)
				break
			}
			if j, err := strconv.ParseFloat(t.Price, 64); err == nil {
				fmt.Println(j)
			}
			s[t.Market] = append(s[t.Market], t.Price)

		}
	}
	for key, element := range s {

		fmt.Println("Market: ", key, "Average Price", element)
	}
}
