package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

var csvFile string

// for each question
var timerDuration int

func init() {
	flag.StringVar(&csvFile, "csv", "problems.csv", "csv file containing math problems")
	flag.IntVar(&timerDuration, "timer", 2, "timer duration for each question")
	flag.Parse()
}

func readCsvRecord(recordCh chan<- []string) {
	defer close(recordCh)
	csvfile, err := os.Open(csvFile)
	if err != nil {
		log.Fatal(err)
	}
	defer csvfile.Close()
	reader := csv.NewReader(csvfile)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		recordCh <- record
	}
}

func main() {
	recordCh := make(chan []string)
	answerCh := make(chan string)
	var correct, incorrect int = 0, 0
	go readCsvRecord(recordCh)

LOOP:
	for record := range recordCh {
		timer := time.NewTimer(time.Duration(timerDuration) * time.Second)
		go func(q string, answerCh chan<- string) {
			fmt.Printf("%s=", q)
			var ans string
			fmt.Scanf("%s", &ans)
			answerCh <- ans
		}(record[0], answerCh)

		select {
		case <-timer.C:
			fmt.Println("timeout!")
			break LOOP
		case userans := <-answerCh:
			// only timer.Stop() doesn't close the timer channel
			// so we have to drain the channel
			if !timer.Stop() {
				<-timer.C
			}
			if userans != record[1] {
				fmt.Println("incorrect")
				incorrect++
			} else {
				fmt.Println("correct")
				correct++
			}
		}
	}
	fmt.Println("Correct answers: ", correct)
	fmt.Println("Incorrect answers: ", incorrect)
}
