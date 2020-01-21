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
	defer file.Close()
	if err != nil {
		log.Print(err)
		return "failed"
	}
	defer file.Close()
	lines, err := ioutil.ReadAll(file)
	if err != nil {
		log.Print(err)
		return "failed"
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
	right := 0
	wrong := 0
	fmt.Printf("WARNING - %v seconds is limit for this test\n", timeout)
	timer := time.NewTimer(time.Duration(timeout) * time.Second)
	for i := 0; i < len(testPoints); i++ {
		fmt.Printf("The question number %v is %v?\nYour answer:", i+1, testPoints[i].question)
		answerChannel := make(chan string)
		go func() {
			var answer string
			scanner := bufio.NewReader(os.Stdin)
			answer, err := scanner.ReadString('\n')
			if err != nil {
				log.Fatal(err)
			}
			answer = strings.Replace(answer, "\n", "", 1)
			answer = strings.Trim(answer, " ")
			answer = strings.ToLower(answer)
			answerChannel <- answer
		}()
		select {
		case <-timer.C:
			fmt.Printf("\nYou are run out of time - %v seconds\n", timeout)
			fmt.Printf("Total number of questions - %v\nThe nubmer of right answers - %v\n", right+wrong, right)
			fmt.Printf("The number of wrong - %v\nPercentage of right - %.1f\n", wrong, float32(right)/float32(right+wrong)*100)
			return
		case answer := <-answerChannel:
			if answer == testPoints[i].answer {
				right++
				fmt.Println("Right one!")
			} else {
				wrong++
				fmt.Println("Wrong one!")
			}
		}
		log.Println("Successfully calculated answers results")
	}
}

func flagsParse() (filePath string, timeout int, shuffle bool) {
	flag.StringVar(&filePath, "path", "./problems.csv", "absolute path to the file")
	flag.IntVar(&timeout, "timeout", 30, "timeout for test")
	flag.BoolVar(&shuffle, "shuffle", false, "shuffle or not")
	flag.Parse()
	return filePath, timeout, shuffle
}

func main() {
	filePath, timeout, shuffle := flagsParse()
	linesFromFile := readLinesFromFile(filePath)
	if linesFromFile == "failed" {
		os.Exit(1)
	}
	linesParsed := parseCSV(linesFromFile)
	var testPoints []testPoint
	testPoints = convertCsvToStructs(linesParsed, shuffle)
	answerAndCalculateResult(testPoints, timeout)
}
