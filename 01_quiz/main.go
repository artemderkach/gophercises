package main

import (
	"encoding/csv"
	"fmt"
	"os"
)

func main() {
	f, err := os.Open("problems.csv")
	if err != nil {
		panic(err)
	}

	csv.NewReader(f)
	fmt.Println("some")
}
