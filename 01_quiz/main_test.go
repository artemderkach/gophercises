package main

import (
	"encoding/csv"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

var problems = `5+5,10
1+1,2
8+3,11
`

// TestQuiz tests without timer component
func TestQuiz(t *testing.T) {
	t.Run("general case", func(t *testing.T) {
		r := csv.NewReader(strings.NewReader(problems))

		result := &Result{}

		records, err := r.ReadAll()
		require.Nil(t, err)

		err = readQuestions(records, getScanner([]string{"10", "4", "11"}), result)
		require.Nil(t, err)
		require.Equal(t, 2, result.CorrectAnswersCount)
	})
	t.Run("test trimming of input", func(t *testing.T) {
		r := csv.NewReader(strings.NewReader(problems))

		records, err := r.ReadAll()
		require.Nil(t, err)

		result := &Result{}
		err = readQuestions(records, getScanner([]string{"10    ", "   4", "       11   "}), result)
		require.Nil(t, err)
		require.Equal(t, 2, result.CorrectAnswersCount)
	})
}

func getScanner(answers []string) func(a ...any) (n int, err error) {
	var i = 0

	// func will iterate through answers
	// length should be controlled outside
	return func(a ...any) (n int, err error) {
		// transform any to string then to string ref
		*a[0].(*string) = answers[i]

		i += 1

		return 0, nil
	}
}
