package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"
)

type Transaction struct {
	transactionID, status string
	timestamp             time.Time
	amount                float64
}

var totalApproved = 0.0
var totalDenied = 0

func worker(roller chan []string, start time.Time, end time.Time) {
	for line := range roller {
		recordCoverted, parseError := time.Parse(time.RFC3339, line[0])
		if parseError != nil {
			fmt.Println("The date is invalid", parseError)
			continue
		}
		if !recordCoverted.Before(start) && !recordCoverted.After(end) {
			f, err := strconv.ParseFloat(line[2], 64)
			if err != nil {
				fmt.Println("Invalid number.", err)
				continue
			}
			tx := Transaction{
				timestamp:     recordCoverted,
				transactionID: line[1],
				amount:        f,
				status:        line[3],
			}

			if tx.status == "approved" {
				totalApproved += tx.amount
			} else {
				totalDenied++
			}
		}
	}
	fmt.Println("Total Approved: ", totalApproved)
	fmt.Println("Total Denied: ", totalDenied)
}

func main() {
	fileFlag := flag.String("file", "", "Choose your file")
	fromDateFlag := flag.String("from", "", "Enter a valid start date")
	toDateFlag := flag.String("to", "", "Enter an end date")

	flag.Parse()

	if *fileFlag == "" {
		fmt.Println("Please, enter a valid file")
		return
	}

	if *fromDateFlag == "" {
		fmt.Println("A start date is missing")
		return
	}

	if *toDateFlag == "" {
		fmt.Println("Final date missing")
		return
	}

	firstDateConverted, errFromDate := time.Parse(time.DateOnly, *fromDateFlag)
	if errFromDate != nil {
		fmt.Println("Enter a valid date", errFromDate)
		return
	}
	finalDateConverted, errToDate := time.Parse(time.DateOnly, *toDateFlag)
	if errToDate != nil {
		fmt.Println("Enter a valid end date", errToDate)
		return
	}

	finalDateConverted = finalDateConverted.AddDate(0, 0, 1)

	if firstDateConverted.After(finalDateConverted) {
		fmt.Println("The interval of dates are invalid.")
		return
	}

	fmt.Println("Chosen file:", *fileFlag)
	fmt.Println("From date:", *fromDateFlag)
	fmt.Println("To date:", *toDateFlag)

	file, err := os.Open(*fileFlag)
	if err != nil {
		fmt.Println("An error occurred", err)
		return
	}

	defer file.Close()

	myReader := csv.NewReader(file)
	myReader.Comma = ';'

	myReader.Read()

	c := make(chan []string)
	go worker(c, firstDateConverted, finalDateConverted)

	for {
		record, errorRead := myReader.Read()
		if errorRead == io.EOF {
			break
		}
		if errorRead != nil {
			fmt.Println("An error had occurred.", errorRead)
			return
		}

		c <- record
	}
	close(c)
}
