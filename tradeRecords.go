package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

type trade struct {
	id     int64   `json:"id"`
	market int64   `json:"market"`
	price  float64 `json:"price"`
	volume float64 `json:volume`
	is_buy bool    `json:"is_buy"`
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	var t trade

	for scanner.Scan() {
		if scanner.Text() == "BEGIN" {
			fmt.Println("Market Open:")
		}
		fmt.Println(scanner.Text())
		x := json.Unmarshal([]byte(scanner.Text()), &t)
		// new func for routune?
		fmt.Println(t.id, t.market, x)

	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
}
