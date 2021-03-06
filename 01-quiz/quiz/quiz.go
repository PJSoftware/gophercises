package quiz

import (
	"encoding/csv"
	"fmt"
	"io"
	"math/rand"
	"os"
	"time"

	"../question"
)

// Quiz holds the entire list of Questions
type Quiz struct {
	questions []*question.Question
	correct   int
	asked     int
}

// Import reads specified CSV file and sets up the quiz
func (qz *Quiz) Import(fileName string, shuffle bool) {
	csvFile, err := os.Open(fileName)
	if err != nil {
		exit(fmt.Sprintf("Unable to open '%s'\n(Error: %v)", fileName, err))
	}
	defer csvFile.Close()

	r := csv.NewReader(csvFile)
	for {
		row, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			exit(fmt.Sprintf("Error reading CSV file: %v", err))
		}
		qz.addQuestion(row[0], row[1])
	}

	if shuffle {
		qz.shuffleQuestions()
	}
}

// Play asks each question in turn
func (qz *Quiz) Play(timeLimit int) {
	nq := len(qz.questions)
	fmt.Printf("Please answer the following %d questions:\n", nq)
	if timeLimit > 0 {
		if timeLimit < nq {
			timeLimit = nq
		}
		fmt.Printf("(You have %d seconds to finish!)\n", timeLimit)
	}

	timer := time.NewTimer(time.Duration(timeLimit) * time.Second)

QuizLoop:
	for _, q := range qz.questions {
		qz.asked++
		fmt.Printf("%d: ", qz.asked)
		go q.Ask()
		select {
		case score := <-q.ChScore:
			qz.correct += score
		case <-timer.C:
			fmt.Println("\nSorry, your time has run out!")
			break QuizLoop
		}
	}
}

// Score displaye the results for the quiz
func (qz *Quiz) Score() {
	nq := len(qz.questions)
	fmt.Printf("You scored %d out of %d\n", qz.correct, nq)
	if qz.correct == nq {
		fmt.Println("Congratulations! You scored 100% correct!")
		return
	}
	if qz.asked < nq {
		fmt.Printf("You only answered %d questions!\n", qz.asked)
		if qz.correct == qz.asked {
			fmt.Println("Of the ones you were asked, you got all correct!")
			return
		}
	}
	fmt.Println("These are the correct answers for the ones you got wrong:")
	for _, q := range qz.questions {
		q.ShowCorrect()
	}
}

func (qz *Quiz) addQuestion(q, a string) {
	qz.questions = append(qz.questions, question.NewQuestion(q, a))
}

func (qz *Quiz) shuffleQuestions() {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(qz.questions), func(i, j int) {
		qz.questions[i], qz.questions[j] = qz.questions[j], qz.questions[i]
	})
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
