package main

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"strconv"
)


// Struct to store raw data
type RawData struct {
	HomeTeam string
	HomeGoals int
	AwayGoals int
	AwayTeam string
}


// Read CSV file having columns "HomeTeam, HomeGoals, AwayGoals, AwayTeam" in that order
func readRawRecordsFromCsv(filepath string) []RawData {
	csvfile, err := os.Open(filepath)
	if err != nil {
		log.Fatalln("Couldn't open the CSV file", err)
	}
	r := csv.NewReader(csvfile)
	lineCount := 0
	var records []RawData
	for {
		lineCount ++
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		if lineCount != 1 {
			homeGoals, strConvErrHome := strconv.Atoi(record[1])
			awayGoals, strConvErrAway := strconv.Atoi(record[2])
			if strConvErrHome != nil {
				log.Fatalln("Error while converting HomeGoals to int", strConvErrHome)
			}
			if strConvErrAway != nil {
				log.Fatalln("Error while converting AwayGoals to int", strConvErrAway)
			}
			record := RawData{
				HomeTeam: record[0],
				HomeGoals: homeGoals,
				AwayGoals: awayGoals,
				AwayTeam: record[3],
			}
			records = append(records, record)
		}
	}
	return records
}