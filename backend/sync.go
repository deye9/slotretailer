package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

// APIlinks is a collection of all needed API links
var APIlinks = make(map[string]string)

// Sync will setup the cadence for sync btw the store and the server.
func Sync() {

	// Get the store details
	GetStore()

	// rotateLogs
	rotateLogs()

	// APIlinks["orders"] = LocalStore.OrdersAPI
	APIlinks["banks"] = LocalStore.BanksAPI
	APIlinks["products"] = LocalStore.ProductsAPI
	APIlinks["customers"] = LocalStore.CustomersAPI

	duration := LocalStore.SyncInterval
	if duration == 0 {
		duration = 5
	}

	tick := time.NewTicker(time.Minute * time.Duration(duration))
	done := make(chan bool)
	go scheduler(tick, done)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
	done <- true
	//os.Exit(1)
}

func scheduler(tick *time.Ticker, done chan bool) {
	task(time.Now())
	for {
		select {
		case t := <-tick.C:
			task(t)
		case <-done:
			return
		}
	}
}

func task(t time.Time) {
	str := strings.ReplaceAll(t.String(), "-", "_")
	str = strings.ReplaceAll(str, ":", "")
	str = strings.Split(str, " ")[0]

	for key, link := range APIlinks {
		go func() {
			// Write the sync start details to the File System.
			WriteFile(BasePath()+"/build/sync/"+str+".log", []byte("Sync for "+key+" started at "+t.String()+"\n"))
		}()

		getAllData(key, link, str)

		time.Sleep(2 * time.Second)
	}
}

func getAllData(key, link, str string) error {
	cmd := ""
	data, err := httpget(link)
	var response interface{}

	if strings.ToLower(key) == "orders" {
		response = []Orders{}
	} else if strings.ToLower(key) == "products" {
		response = []Products{}
	} else if strings.ToLower(key) == "customers" {
		response = []Customers{}
	} else if strings.ToLower(key) == "banks" {
		response = []Banks{}
	}

	if err = json.Unmarshal(data, &response); err != nil {
		CheckError("Error unmarshalling data.", err, false)
		return err
	}
	cmd = structToInsertUpdate(response, key)

	if err = Modify(cmd); err != nil {
		CheckError("Error saving HTTPGET result. ", err, false)
		return err
	}

	go func() {
		// Write the sync start details to the File System.
		WriteFile(BasePath()+"/build/sync/"+str+".log", []byte("Sync for "+key+" finished at "+time.Now().String()+"\n"))
	}()
	return nil
}

func httpget(url string) (data []byte, err error) {
	var res *http.Response

	// make a HTTP GET request
	if res, err = http.Get(url); err != nil {
		CheckError("Error getting data at URL: "+url, err, false)
		return
	}

	// read all response body
	if data, err = ioutil.ReadAll(res.Body); err != nil {
		CheckError("Error reading GET response body. ", err, false)
	}

	// close response body
	res.Body.Close()

	return
	// // print `data` as a string
	// fmt.Printf("%s\n", data)
}

func httppost(url, payload string) (data []byte, err error) {
	// Set the requestBody payload
	var res *http.Response
	requestBody := strings.NewReader(payload)

	// post the data to the server
	if res, err = http.Post(url, "application/json; charset=UTF-8", requestBody); err != nil {
		CheckError("Error POSTING data to URL: "+url, err, false)
		return
	}

	// read response data
	data, _ = ioutil.ReadAll(res.Body)

	// close response body
	res.Body.Close()

	return
	// // print request `Content-Type` header
	// requestContentType := res.Request.Header.Get("Content-Type")
	// fmt.Println("Request content-type:", requestContentType)
}

// GetLogs returns a list of all available log files.
func GetLogs() (files []string, err error) {
	// read the /build/sync/ directory and returns a list of files sorted by filename.
	fileInfo, err := ioutil.ReadDir(BasePath() + "/build/sync/")
	if err != nil {
		return nil, err
	}

	for _, file := range fileInfo {
		files = append(files, strings.Replace(file.Name(), ".log", "", -1))
	}

	return files, nil
}

// GetLog returns the details of the log file requested.
func GetLog(id string) (fileContent string, err error) {
	var fileInfo []byte

	if fileInfo, err = ReadFile(BasePath() + "/build/sync/" + id); err != nil {
		return "", err
	}

	return string(fileInfo), nil
}

// rotateLogs rotates the sync logs based on the configured Frequency
func rotateLogs() {

	// Get the default time Calculations
	t := time.Now()
	dateRange := []string{}
	year, currentWeek := t.ISOWeek()
	yesterday := t.AddDate(0, 0, -1)
	lastMonth := t.AddDate(0, -1, 0)
	lastQuarter := t.AddDate(0, -3, 0)
	lastYear := t.AddDate(-1, 0, 0)

	switch strings.ToLower(LocalStore.LogRotation) {
	case "daily":
		dateRange = append(dateRange, yesterday.String())

	case "weekly":
		start, end := WeekRange(year, currentWeek-1)
		dateRange = append(dateRange, start.String())
		for i := 1; i <= 5; i++ {
			dateRange = append(dateRange, start.AddDate(0, 0, i).String())
		}
		dateRange = append(dateRange, end.String())

	case "bi-weekly":
		start := WeekStart(year, currentWeek-2)
		dateRange = append(dateRange, start.String())
		for i := 1; i <= 13; i++ {
			dateRange = append(dateRange, start.AddDate(0, 0, i).String())
		}

	case "monthly":
		dateRange = append(dateRange, lastMonth.String())

	case "quarterly":
		dateRange = append(dateRange, lastQuarter.String())

	case "yearly":
		dateRange = append(dateRange, lastYear.String())
	}

	for _, date := range dateRange {
		str := strings.ReplaceAll(date, "-", "_")
		str = strings.ReplaceAll(str, ":", "")
		str = strings.Split(str, " ")[0] + ".log"
		deleteFile(BasePath() + "/build/sync/" + str)
	}
}
