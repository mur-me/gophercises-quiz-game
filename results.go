package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
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

func convertCsvToStructs(line [][]string, shuffle bool) []testPoint {
	var testPoints []testPoint
	for i := 0; i < len(line); i++ {
		question := line[i][0]
		answer := line[i][1]
		answer = strings.Trim(answer, " ")
		answer = strings.ToLower(answer)
		testPoint := testPoint{question: question, answer: answer}
		testPoints = append(testPoints, testPoint)
	}
	log.Println("Converted CSV to slices of structs")
	if shuffle {
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(testPoints), func(i, j int) { testPoints[i], testPoints[j] = testPoints[j], testPoints[i] })
		log.Println("Shuffled slice of questions")
	}
	return testPoints
}

func answerAndCalculateResult(testPoints []testPoint, timeout int) {
	var right, wrong int
	fmt.Printf("WARNING - %v seconds is limit for this test\n", timeout)
	channelForTimeoutCheck := make(chan string, 1)
	go func() {
		for i := 0; i < len(testPoints); i++ {
			scanner := bufio.NewReader(os.Stdin)
			fmt.Printf("The question number %v is %v?\nYour answer:", i+1, testPoints[i].question)
			answer, err := scanner.ReadString('\n')
			if err != nil {
				log.Fatal(err)
			}
			answer = strings.Replace(answer, "\n", "", 1)
			answer = strings.Trim(answer, " ")
			answer = strings.ToLower(answer)
			if answer != testPoints[i].answer {
				wrong++
				fmt.Println("Wrong one!")
			} else {
				right++
				fmt.Println("Right one!")
			}
		}
		var rightPerc float32
		rightPerc = float32(right) / float32(right+wrong) * 100
		totalQuestions := right + wrong

		fmt.Printf("Total number of questions - %v\nThe nubmer of right answers - %v\n", totalQuestions, right)
		fmt.Printf("The number of wrong - %v\nPercentage of right - %.1f\n", wrong, rightPerc)
		log.Println("Successfully calculated answers results")
		channelForTimeoutCheck <- "Hooray, u are finished in time"
	}()
	select {
	case result := <-channelForTimeoutCheck:
		fmt.Println(result)
	case <-time.After(time.Duration(timeout) * time.Second):
		fmt.Println()
		log.Fatalf("Program finished after timeout in %v seconds", timeout)
	}
}

func flagsParse() (string, int, bool) {
	var filePath string
	var timeout int
	var shuffle bool
	flag.StringVar(&filePath, "path", "./problems.csv", "absolute path to the file")
	flag.IntVar(&timeout, "timeout", 30, "timeout for test")
	flag.BoolVar(&shuffle, "shuffle", false, "shuffle or not")
	flag.Parse()
	return filePath, timeout, shuffle
}

func main() {
	filePath, timeout, shuffle := flagsParse()
	linesFromFile := readLinesFromFile(filePath)
	linesParsed := parseCSV(linesFromFile)
	var testPoints []testPoint
	testPoints = convertCsvToStructs(linesParsed, shuffle)
	answerAndCalculateResult(testPoints, timeout)
}
