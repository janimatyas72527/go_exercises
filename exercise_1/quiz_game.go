package exercise_1

import (
	"bufio"
	"context"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"time"
)

func QuizGame() {
	// Command-line parameters
	file_ptr := flag.String("file", "./problems.csv", "File in which the questions are found")
	duration_ptr := flag.Duration("timeout", 30*time.Second, "The timeout in seconds")
	use_time_ptr := flag.Bool("timedquiz", false, "Use timed quiz")
	shuffle_ptr := flag.Bool("shuffle", false, "Shuffle questions")
	flag.Parse()

	file, err_file := filepath.Abs(*file_ptr)
	if err_file != nil {
		log.Fatal(err_file)
	}

	// Open file
	f, err_f := os.Open(file)
	if err_f != nil {
		log.Fatal(err_f)
	}

	fmt.Println("File: ", file)
	fmt.Println("Use Timed Quiz: ", *use_time_ptr)
	fmt.Println("Time: ", *duration_ptr)
	fmt.Println("Shuffle Questions: ", *shuffle_ptr)

	// Read csv values using csv.Reader
	csvReader := csv.NewReader(f)
	questions, err_f := csvReader.ReadAll()
	if err_f != nil {
		log.Fatal(err_f)
	}
	f.Close()
	// Shuffle questions if requested
	if *shuffle_ptr {
		swap := reflect.Swapper(questions)
		rand.Shuffle(len(questions), func(a int, b int) { swap(a, b) })
	}

	// Initialize
	reader := bufio.NewReader(os.Stdin)
	answers := make([]string, len(questions))
	correct_answers := 0

	// Lambda for asking questions
	ask_questions := func(questions [][]string, answers []string, reader *bufio.Reader, ctx context.Context) {
		for index := 0; index < len(questions); index++ {
			before_deadline := true
			if ctx != nil {
				deadline, _ := ctx.Deadline()
				before_deadline = time.Now().Before(deadline)

			}
			if before_deadline {
				question := questions[index][0]
				fmt.Print("Question: ", question, "? Answer: ")
				answer, err_answer := reader.ReadString('\n')
				if err_answer != nil {
					log.Fatal(err_answer)
				} else {
					trimmed_answer := ""
					if runtime.GOOS == "windows" {
						trimmed_answer = strings.TrimRight(answer, "\r\n")
					} else {
						trimmed_answer = strings.TrimRight(answer, "\n")
					}
					answers[index] = trimmed_answer
				}
			} else {
				break
			}
		}
	}

	// Two ways of quiz handling
	if *use_time_ptr {
		// Timed Quiz
		fmt.Println("Please hit 'Enter' to start the timed quiz.")
		reader.ReadString('\n')
		deadline_time := time.Now().Add(*duration_ptr)
		ctx, cf := context.WithDeadline(context.Background(), deadline_time)
		go ask_questions(questions, answers, reader, ctx)
		time.Sleep(*duration_ptr)
		cf()
	} else {
		// Untimed Quiz
		ask_questions(questions, answers, reader, nil)
	}

	// Evaluate answers
	for index := 0; index < len(questions); index++ {
		if answers[index] == questions[index][1] {
			correct_answers++
		} else {
			fmt.Println("For question: '", questions[index][0], "' your answer is: '",
				answers[index], "' but correct answer is: '", questions[index][1], "'.")
		}
	}

	// Print results
	fmt.Println("You have ", correct_answers, " correct answers out of ", len(questions), ".")
}
