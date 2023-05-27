package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

// Result stores the data about quiz state for questions and answers
type Result struct {
	QuestionsCount      int
	CorrectAnswersCount int
}

// Scanner is an interface to repeat fmt.Scan function
type Scanner func(a ...any) (n int, err error)

func main() {
	var fname string
	var timer int
	flag.StringVar(&fname, "f", "problems.csv", "filename")
	flag.IntVar(&timer, "t", 30, "time limit")
	flag.Parse()

	f, err := os.Open(fname)
	if err != nil {
		log.Fatalln(err)
	}

	r := csv.NewReader(f)

	result := &Result{}

	fmt.Println("running file:", fname)
	fmt.Println("time limit:", timer)
	fmt.Print("Press any Enter to start the test")
	_, err = fmt.Scanln()
	if err != nil {
		log.Fatalln(err)
	}

	// put quiz in separate routine to run it on timer
	go func() {
		err = quiz(r, fmt.Scan, result)
		if err != nil {
			log.Fatalln(err)
		}
	}()

	<-time.NewTimer(time.Second * time.Duration(timer)).C

	fmt.Println()
	fmt.Println("score:", result.CorrectAnswersCount, "/", result.QuestionsCount)
}

// run quiz
func quiz(r *csv.Reader, f Scanner, result *Result) error {
	records, err := r.ReadAll()
	if err != nil {
		return err
	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(records), func(i, j int) { records[i], records[j] = records[j], records[i] })

	return readQuestions(records, f, result)
}

// start reading question one by one
func readQuestions(records [][]string, f Scanner, result *Result) error {
	result.QuestionsCount = len(records)

	// record[0] - question
	// record[1] - answer
	for _, record := range records {
		fmt.Print(record[0], "? ")

		// read user input
		var input string
		_, err := f(&input)
		if err != nil {
			return err
		}

		// compare user input with correct answer
		if record[1] == strings.TrimSpace(input) {
			result.CorrectAnswersCount += 1
		}
	}

	return nil
}
