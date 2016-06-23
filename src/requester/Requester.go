package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// TokenAPI BSW TokenAPI
const TokenAPI = "https://api.bidswitch.com/discrepancy-check/v1.0/login"

// API BSW post api
const API = "https://api.bidswitch.com/discrepancy-check/v1.0/ssp/imobile/upload-report/"

// User get the stupid token
type User struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

// TokenResponse Unmarshal the stupid TokenResponse
type TokenResponse struct {
	AccessToken string `json:"access_token"`
}

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

func getToken() (tokenStr string) {
	tokenStr = ""
	user := User{UserName: "*************", Password: "**************"}
	buffer := new(bytes.Buffer)
	json.NewEncoder(buffer).Encode(user)

	tokenRes, _ := http.Post(TokenAPI, "application/json; charset=utf-8", buffer)

	if tokenRes.StatusCode == 200 {
		bodyBytes, _ := ioutil.ReadAll(tokenRes.Body)
		token := new(TokenResponse)
		json.Unmarshal(bodyBytes, &token)
		tokenStr = token.AccessToken
	}
	return
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

	fmt.Println(jsonstring)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Get the stupid token
	token := getToken()

	tokenFile, _ := os.Create("../datasource/token.txt")

	tokenFile.Write([]byte(token))
	tokenFile.Close()

	// Use the stupid token to post... Why don't use basic auth directly stupid BSW
	req, _ := http.NewRequest("POST", API, strings.NewReader(jsonstring))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", token)

	client := &http.Client{Timeout: time.Duration(15 * time.Second)}
	res, apierr := client.Do(req)
	defer res.Body.Close()

	if apierr != nil {
		fmt.Println(apierr)
	} else {
		fmt.Println(res)
	}
}
