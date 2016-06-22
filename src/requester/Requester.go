package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// API BSW post api
const API = "https://api.bidswitch.com/discrepancy-check/v1.0/ssp/imobile/upload-report/"

// DailyData atomic daily data
type DailyData struct {
	Cost float64 `json:"cost"`
	Imps int     `json:"imps"`
}

// Data data in each record
type Data struct {
	DailyData []DailyData `json:"$date$"`
}

// Record each line of csv
type Record struct {
	SeatID   string `json:"seat"`
	Currency string `json:"currency"`
	Data     Data   `json:"data"`
	TimeZone string `json:"timezone"`
}

// Report the report data
type Report struct {
	Report []Record `json:"reports"`
}

var reportdate string

func resloveDataSource() string {
	return "../datasource/report_" + reportdate + ".csv"
}

func init() {
	flag.StringVar(&reportdate, "date", time.Now().Format("2006-01-02"), "date of the report")
}

func main() {
	csvFile, err := os.Open(resloveDataSource())

	if err != nil {
		fmt.Println(err)
	}

	defer csvFile.Close()

	reader := csv.NewReader(csvFile)

	reader.FieldsPerRecord = 4

	csvData, err := reader.ReadAll()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var oneRecord Record
	var allRecords []Record

	for index, row := range csvData {
		// skip header
		if index == 0 {
			continue
		}

		seatID := row[0]
		currency := row[1]
		imps, _ := strconv.Atoi(row[2])
		cost, _ := strconv.ParseFloat(row[2], 64)
		timezone := "Asia/Tokyo" // fixed

		oneRecord = Record{SeatID: seatID, Currency: currency, Data: Data{DailyData: []DailyData{{Cost: cost, Imps: imps}}}, TimeZone: timezone}
		allRecords = append(allRecords, oneRecord)
	}

	report := Report{Report: allRecords}

	jsondata, err := json.Marshal(report)

	// for the sake of BSW stupid api
	jsonstring := strings.Replace(string(jsondata[:]), "$date$", reportdate, -1)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	jsonFile, err := os.Create("../datasource/data.json")

	if err != nil {
		fmt.Println(err)
	}

	defer jsonFile.Close()

	jsonFile.Write([]byte(jsonstring))
	jsonFile.Close()
}
