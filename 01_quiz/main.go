package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

var counter = 1
var correctCounter = 0

func main() {
	timeLimit := flag.Int64("t", 10, "time limit")
	inputFile := flag.String("f", "problems.csv", "file with problems")
	flag.Parse()

	f, err := os.Open(*inputFile)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = f.Close()
	}()

	go func() {
		err = startQuiz(f)
		if err != nil {
			fmt.Println("error:", err)
		}
		os.Exit(0)
	}()

	t := time.NewTimer(time.Second * time.Duration(*timeLimit))
	<-t.C

	fmt.Printf("\nresult: %d/%d\n", correctCounter, counter)
}

func startQuiz(f io.Reader) error {
	r := csv.NewReader(f)
	// read headers
	_, err := r.Read()
	if err != nil {
		return err
	}

	for record, err := r.Read(); ; record, err = r.Read() {
		// check for unexpected errors
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		problem, solution := record[1], record[2]

		fmt.Printf("Problem #%d: %s? ", counter, problem)

		// get the input from command line
		var input string
		_, err := fmt.Scan(&input)
		if err != nil {
			return err
		}
		input = strings.Trim(input, "\n")

		counter += 1
		// check the correctness of answer
		if solution == input {
			correctCounter += 1
		}
	}

	fmt.Printf("result: %d/%d\n", correctCounter, counter)

	return nil
}
