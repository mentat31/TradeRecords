package funcs

//	Functions and types for aggregating trade data.

import (
	"encoding/json"
	"math"

	"github.com/montanaflynn/stats"
)

/*
	Trade represents the current trade being read in. Struct fields are set to match
	input, and exported for Json recognition.
*/

type Trade struct {
	Id     json.Number `json:"id"`
	Market json.Number `json:"market" type:"integer" required:"true"`
	Price  json.Number `json:"price" type:"float64" required:"true"`
	Volume json.Number `json:"volume"`
	Is_buy bool        `json:"is_buy"`
}

/*
	Persist represents the data in between read and write. It serves as the value
	in RecMap (see below) and is indexed by maket. The fields Prices and Volume append
	each data point in a slice when read. Percentage_Buy stores a float64(1) on true.
	This slice serves as the record of buys and is used in calculating the proportion.
*/

type Persist struct {
	Volume         []float64
	Prices         []float64
	Percentage_buy []float64
	VWAP           float64
}

/*
	Record represents the final output before Json marshalling. Struct fields contain
	output of Agg function.
*/

type Record struct {
	Market         json.Number `json:"market"`
	TotalVolume    float64     `json:"Total Volume"`
	MeanPrice      float64     `json:"Mean Price"`
	MeanVolume     float64     `json:"Mean Volume"`
	VWAP           float64     `json:"VWAP"`
	Percentage_buy float64     `json:"Percentage Buy"`
}

// RecMap is indexed by Market and holds its aggregate data.

type RecMap map[json.Number]Persist

func Agg(current Trade, record Persist, marketData RecMap) Persist {

	/*
		Side Note 1:

		Coming from the world of Scala, I was really looking forward to exploring
		concurrency in Go. However, I got some very explicit error messages while doing so.
		Not being too familiar with the language, It's definitely something I am looking
		into, however it was my initial approach.

		Side Note 2:

		During my initial tests, I noticed time to complete had slightly increased by about
		10ms when importing the functions and structs into the main file. I know the binary
		is not deterministic, however I wanted to keep the types and functions organized in
		a single module.
	*/

	// Get/set current market data in RecMap.
	r := marketData[current.Market]

	// Convert Price and Volume from Json.Number to float64 for operations.
	p, _ := current.Price.Float64()
	v, _ := current.Volume.Float64()

	// Append trade's price and volume to Market's Persist struct.
	r.Prices = append(r.Prices, p)
	r.Volume = append(r.Volume, v)

	/*
		Volume Weighted Average Price accumulates at each trade.
		Specifically "Typical Price" ( (High + Low + Close {Last in this instance} ) / 3) * Volume.
		The final output is then divided by the total volume before output.

	*/
	max, _ := stats.Max(r.Prices)
	min, _ := stats.Min(r.Prices)
	r.VWAP += ((max + min + p) / 3) * v

	// A count of every buy. Divided by total trades before output.

	if current.Is_buy == true {
		r.Percentage_buy = append(r.Percentage_buy, float64(1))
	}
	// Market data set to Market index in RecMap
	marketData[current.Market] = r
	return marketData[current.Market]
}

// A rounding function, courtesy of https://yourbasic.org/golang/round-float-2-decimal-places/

func RoundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}
