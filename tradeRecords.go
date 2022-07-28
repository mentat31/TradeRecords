package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type trade struct {
	Id     json.Number `json:"id"`
	Market json.Number `json:"market"`
	Price  json.Number `json:"price"`
	Volume json.Number `json:"volume"`
	Is_buy bool        `json:"is_buy"`
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	//var t trade

	for scanner.Scan() {
		var t trade

		switch scanner.Text() {

		case "BEGIN":
			fmt.Println("Market Open:")
		case "END":
			break
		default:
			x := json.NewDecoder(strings.NewReader(scanner.Text())).Decode(&t)

			// new func for routune?
			if x != nil {
				fmt.Println(scanner.Text(), "Wont work because", x)
				break

			} else {
				fmt.Println(t.Id)
				break
			}

		}
	}

}
