package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

const problemsFileName = "problems.csv"

var (
	correctAnswers, totalQuestions int
)

func main() {
	// flag package used to get arguments from user in runtime in cmd
	var (
		flagFileName = flag.String("p", problemsFileName, "The path to the problems csv file")
		flagTimer    = flag.Duration("t", 30*time.Second, "The Max time for quiz")
	)

	// start reading/parsing the above defined flags
	flag.Parse()

	if flagFileName == nil || flagTimer == nil {
		fmt.Println("Missing arguments of filename and/or timer")
		return
	}

	fmt.Printf("Hit enter to start quiz form %q in %v\n", *flagFileName, *flagTimer)

	fmt.Scanln() // wait for enter key

	// x := 23 // to define varaibles into function scope (automatically go define the data type)
	// var x int = 23 // explicitly define the datatype to int

	// read csv file (problems.csv) using package
	// flagFileName is pointer and this method take string not point, we use add * before it to convert it to normal string
	// * is inverse of & (return/convert the pointer)
	file, err := os.Open(*flagFileName)
	if err != nil {
		fmt.Printf("Failed to load file %v", err)
		return
	}
	defer file.Close() // run this line after the main function is done executing

	r := csv.NewReader(file)

	questions, err := r.ReadAll()
	if err != nil {
		fmt.Printf("Failed to read csv file %v", err)
		return
	}
	totalQuestions := len(questions)

	// start 30 seconds timer
	// will do it into goroutine and will return channel while executing
	doneQuize := StartQuiz(questions)

	// fire a timer in a separate channel
	timerQuize := time.NewTimer(*flagTimer).C

	// select is like switch but for channels
	select {
	case <-doneQuize:
		fmt.Println("all answers are done")
	case <-timerQuize:
		fmt.Println("the time is over")
	}

	// output number of questions (total + correct)
	fmt.Printf("Result: %d/%d \n", correctAnswers, totalQuestions)
}

func StartQuiz(questions [][]string) chan bool {
	// go keyword is use to define goroutine (a lighweight thread to be executed into parallel execution)
	// - must written before function execution
	// chan keyword (refer to channels) used to define channels to communicate between different goroutines
	// Channels are pipes through which you can send and receive values using the channel operator, <- .
	done := make(chan bool)

	go func() {
		for i, record := range questions {
			var (
				question      = record[0]
				correctAnswer = record[1]
			)

			// display one question at a time
			fmt.Printf("%d. %s?\n", i+1, question)
			var userAnswer string

			// get answer from, then proceed to the next question immediately
			_, err := fmt.Scan(&userAnswer) // read input from user
			if err != nil {
				fmt.Printf("Failed to scan: %v", err)
				return
			}

			userAnswer = strings.TrimSpace(userAnswer)
			userAnswer = strings.ToLower(userAnswer)
			if userAnswer == correctAnswer {
				correctAnswers++
			}
		}
		// notify the main thread that we're done runnning the quiz
		done <- true
	}()

	return done
}
