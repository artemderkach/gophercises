package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

var wg = sync.WaitGroup{}

var defaultTimer int = 30
var defautlFile string = "./problems.csv"

func main() {
	quizLimitFlag := flag.Int("limit", defaultTimer, "number seconds that will be used for quiz")
	filePathFlag := flag.String("file", defautlFile, "file with questions")
	isShuffleFlag := flag.Bool("shuffle", false, "shuffle file records")
	flag.Parse()

	filePath, err := filepath.EvalSymlinks(*filePathFlag)
	if err != nil {
		panic(err)
	}
	wg.Add(1)
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	r := csv.NewReader(file)

	go runQuiz(r, *quizLimitFlag, *isShuffleFlag)
	wg.Wait()
}

func runQuiz(r *csv.Reader, seconds int, isShuffle bool) {
	records, err := r.ReadAll()
	if isShuffle {
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(records), func(i, j int) { records[i], records[j] = records[j], records[i] })
	}
	if err != nil {
		panic(err)
	}

	numberOfAnswers := 0
	numberOfQuestions := len(records)
	go func() {
		<-time.NewTimer(time.Duration(seconds) * time.Second).C
		end(numberOfAnswers, numberOfQuestions)
	}()

	for _, record := range records {
		question := record[0]
		answer := record[1]
		fmt.Print(question, ": ")
		reader := bufio.NewReader(os.Stdin)
		userAnswer, _ := reader.ReadString('\n')
		userAnswer = strings.ToLower(strings.Trim(userAnswer[:len(userAnswer)-1], " "))
		if userAnswer == answer {
			numberOfAnswers++
		}
	}
	end(numberOfAnswers, numberOfQuestions)
}

func end(numberOfAnswers int, numberOfQuestions int) {
	result := "Result:" + strconv.Itoa(numberOfAnswers) + "/" + strconv.Itoa(numberOfQuestions)
	fmt.Println("")
	fmt.Println(result)
	fmt.Println("END")
	wg.Done()
}
