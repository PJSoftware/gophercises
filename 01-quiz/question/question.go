package question

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Question is an individual q&a combo
type Question struct {
	question string
	answer   string
	response string
	correct  bool
	asked    bool
}

// NewQuestion constructs a new Question struct
func NewQuestion(question, answer string) *Question {
	q := new(Question)
	q.question = question
	q.answer = answer
	return q
}

// Ask handles asking (and scoring the answer of) each question
func (q *Question) Ask() int {
	score := 0

	fmt.Printf("%s = ? ", q.question)
	q.asked = true
	q.UserInput()

	if q.response == q.answer {
		score = 1
		q.correct = true
	}

	return score
}

// UserInput reads the user's response
func (q *Question) UserInput() {
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	text = strings.Replace(text, "\n", "", -1) // for all platforms
	text = strings.Replace(text, "\r", "", -1) // for Windows users
	q.response = text
}

// ShowCorrect displays the correct answer for any question
// the user got wrong
func (q *Question) ShowCorrect() {
	if q.asked && !q.correct {
		fmt.Printf(" '%s' is '%s'; you said '%s'\n", q.question, q.answer, q.response)
	}
}