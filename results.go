package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func readLinesFromFile(path string) string {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	lines, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("File succesfully parsed")
	return string(lines)
}

func parseCSV(lines string) [][]string {
	reader := csv.NewReader(strings.NewReader(lines))
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("CSV lines succesfully parsed")
	return records
}

type testPoint struct {
	question, answer string
}

func convertCsvToStructs(line [][]string) []testPoint {
	var testPoints []testPoint
	for i := 0; i < len(line); i++ {
		question := line[i][0]
		answer := line[i][1]
		testPoint := testPoint{question: question, answer: answer}
		testPoints = append(testPoints, testPoint)
	}
	log.Println("Converted CSV to slices of structs")
	return testPoints
}

func answerAndCalculateResult(testPoints []testPoint, timeout int) {
	var right, wrong int
	channelForTimeout := make(chan string, 1)
	go func() {
		channelForTimeout <- "finish"
	}()
	for i := 0; i < len(testPoints); i++ {
		scanner := bufio.NewReader(os.Stdin)
		fmt.Printf("The question number %v is %v:", i, testPoints[i].question)
		answer, err := scanner.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		answer = strings.Replace(answer, "\n", "", 1)
		if answer != testPoints[i].answer {
			wrong++
		} else {
			right++
		}
	}
	var rightPerc float32
	rightPerc = float32(right) / float32(right+wrong) * 100
	totalQuestions := right + wrong
	fmt.Printf("Total number of questions - %v\nThe nubmer of right answers - %v\n", totalQuestions, right)
	fmt.Printf("The number of wrong - %v\nPercentage of right - %.1f\n", wrong, rightPerc)
	log.Println("Successfully calculated answers results")
}

func flagsParse() (string, int) {
	var filePath string
	var timeout int
	flag.StringVar(&filePath, "path", "./problems.csv", "absolute path to the file")
	flag.IntVar(&timeout, "timeout", 30, "timeout for test")
	flag.Parse()
	return filePath, timeout
}

func main() {
	filePath, timeout := flagsParse()
	linesFromFile := readLinesFromFile(filePath)
	linesParsed := parseCSV(linesFromFile)
	var testPoints []testPoint
	testPoints = convertCsvToStructs(linesParsed)
	answerAndCalculateResult(testPoints, timeout)
}
