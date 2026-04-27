package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"time"
)

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

	record, readError := myReader.Read()

	if readError != nil {
		fmt.Println("It was impossible to read the file", readError)
		return
	}

	fmt.Println(record)

}
