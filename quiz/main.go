package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

//Problem contains a single Question and Answer
type Problem struct {
	Question string
	Answer   string
}

var score int

func answerRoutine(problems []Problem, expired chan bool, timer *time.Timer, duration int) {
	stdin := bufio.NewReader(os.Stdin)
	for _, problem := range problems {
		fmt.Printf("What is %v? ", problem.Question)
		ans, _ := stdin.ReadString('\n')
		if strings.TrimSpace(ans) == problem.Answer {
			score++
			fmt.Println("Correct!")
		} else {
			fmt.Println("Oops, you didn't answer that one correctly.")
			fmt.Printf("The correct answer is %v.\n", problem.Answer)
		}
		if !timer.Stop() {
			<-timer.C
		}
		timer.Reset(time.Duration(duration) * time.Second)
	}
	expired <- false
}

func timerChecker(done chan bool, timer *time.Timer) {
	<-timer.C
	done <- true
}

func main() {
	var csvFile string
	var timeLimit int
	flag.StringVar(&csvFile, "csvpath", "", "Path to CSV File")
	flag.IntVar(&timeLimit, "timelimit", 30, "Time limit per question")
	flag.Parse()

	csvData, e := os.Open(csvFile)
	if e != nil {
		panic(e)
	}

	defer csvData.Close()

	reader := csv.NewReader(csvData)
	var problems []Problem
	for {
		line, e := reader.Read()
		if e == io.EOF {
			break
		} else if e != nil {
			panic(e)
		}
		problems = append(problems, Problem{
			line[0],
			line[1],
		})
	}
	expired := make(chan bool)
	timer := time.NewTimer(time.Duration(timeLimit) * time.Second)

	go answerRoutine(problems, expired, timer, timeLimit)
	go timerChecker(expired, timer)

	val := <-expired
	fmt.Println("")

	if val {
		fmt.Println("Time's UP!")
	} else {
		fmt.Println("You reached the end of the quiz.")
	}

	fmt.Printf("You answered %d out of %d questions correctly.", score, len(problems))
}
