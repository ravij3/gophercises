/*
 * Copyright (c) 2019, Ravi Jadhav. All rights reserved.
 */

package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"
	"time"
)

const DefaultFileName = "/path/to/file"
const Timer = 10 //seconds

type pair = struct {
	question, answer string
}
type QuestionAnswer = map[int]pair

func main() {
	fileNameFlag := flag.String("f", DefaultFileName, "File name")
	timerFlag := flag.Int("t", Timer, "Total time for quiz")
	shuffleFlag := flag.Bool("s", false, "Shuffle Questions")
	flag.Parse()

	csvFile, err := os.Open(*fileNameFlag)
	if err != nil {
		fmt.Println(err)
		return
	}

	userInput := bufio.NewScanner(os.Stdin)
	csvReader := csv.NewReader(bufio.NewReader(csvFile))
	total, correct := 0, 0
	systemQA := make(QuestionAnswer)
	userQA := make(QuestionAnswer)

	for {
		line, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println(err)
			return
		}
		systemQA[total] = pair{line[0], line[1]}
		total++
	}

	questionNos := make([]int, total-1)
	if !(*shuffleFlag) {
		for i := range questionNos {
			questionNos[i] = i
		}
	} else {
		questionNos = rand.Perm(total)
	}


	fmt.Println("Enter to begin quiz")

	ansChan := make(chan string)
	timer := time.NewTimer(time.Duration(*timerFlag) * time.Second)
	defer timer.Stop()

	userInput.Scan()

outer:
	for _, val := range questionNos {
		fmt.Println(systemQA[val].question)

		go func() {
			userInput.Scan()
			res := userInput.Text()
			ansChan <- res
		}()

		select {
		case <-timer.C:
			break outer
		case ans := <-ansChan:
			userQA[val] = pair{"", ans}
		}
	}

	correct = calculateCorrectAnswers(userQA, systemQA)
	fmt.Printf("Total questions %d\n", total)
	fmt.Printf("Total correct answers %d", correct)
}

func calculateCorrectAnswers(userQA, systemQA map[int]pair) int {
	correct := 0
	for key, val := range userQA {
		if strings.EqualFold(strings.TrimSpace(val.answer), strings.TrimSpace(systemQA[key].answer)) {
			correct++
		}
	}
	return correct
}

