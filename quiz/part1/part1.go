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
	"os"
)

const DefaultFileName = "/Users/ravij/golearning/src/gophercises/quiz/problems.csv"

func main() {
	fileNameFlag := flag.String("f", DefaultFileName, "File name")
	flag.Parse()

	csvFile, err := os.Open(*fileNameFlag)
	if err!=nil {
		fmt.Println(err)
		return
	}

	userInput := bufio.NewScanner(os.Stdin)
	csvReader := csv.NewReader(bufio.NewReader(csvFile))
	total, correct := 0, 0

	for {
		line, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(line[0])

		userInput.Scan()
		res := userInput.Text()
		if res == line[1] {
			correct++
		}
		total++
	}

	fmt.Printf("Total questions %d\n", total)
	fmt.Printf("Total correct answers %d", correct)
}

