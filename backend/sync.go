package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
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

	// rotateLogs if folder exists else create the sync folder
	if _, err := os.Stat(BasePath() + "/build/sync"); os.IsNotExist(err) {
		if err := os.Mkdir(BasePath() + "/build/sync", 0755); err != nil {
			CheckError("Error ", errors.New("creating the Sync folder"), false)
		}
	} else {
		go rotateLogs()
	}

	if LocalStore.BanksAPI == "" || LocalStore.ProductsAPI == "" || LocalStore.CustomersAPI == "" {
		CheckError("LocalStore endpoint missing.", errors.New("Missing endpoint from application"), false)
		return
	}

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
	<-sigs // recieve only channel
	done <- true
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
		// Write the sync start details to the File System via a Goroutine.
		go WriteFile(BasePath()+"/build/sync/"+str+".log", []byte("Sync for "+key+" started at "+t.String()+"\n"))

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

	// Write the sync start details to the File System via a Goroutine.
	go WriteFile(BasePath()+"/build/sync/"+str+".log", []byte("Sync for "+key+" finished at "+time.Now().String()+"\n"))
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

func sapAuthenticate() {
	url := "https://197.255.32.34:50000/b1s/v1/Login"
	method := "POST"

	payload := strings.NewReader("{\n    \"CompanyDB\": \"SlotTESTDB\",\n    \"UserName\": \"Tech1\",\n    \"Password\": \"P@ss0123\"\n}")

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Cookie", "B1SESSION=66f5ed8a-0bb5-11eb-8000-00155d0abe08; ROUTEID=.node7")

	res, err := client.Do(req)
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	fmt.Println(string(body))
}

// rotateLogs rotates the sync logs based on the configured Frequency
func rotateLogs() {
	// Get the default time Calculations
	const (
		month    = 1
		quarter  = 3
		fullYear = 12
	)
	t := time.Now()
	year, currentWeek := t.ISOWeek()
	yesterday := t.AddDate(0, 0, -1)

	switch strings.ToLower(LocalStore.LogRotation) {
	case "daily":
		deleteFile(BasePath() + "/build/sync/" + yesterday.Format("2006_01_02.log"))

	case "weekly":
		start := WeekStart(year, currentWeek)
		deleteFile(BasePath() + "/build/sync/" + start.Format("2006_01_02.log"))
		for i := 1; i <= 5; i++ {
			deleteFile(BasePath() + "/build/sync/" + start.AddDate(0, 0, i).Format("2006_01_02.log"))
		}

	case "bi-weekly":
		start := WeekStart(year, currentWeek-1)
		deleteFile(BasePath() + "/build/sync/" + start.Format("2006_01_02.log"))
		for i := 1; i <= 13; i++ {
			deleteFile(BasePath() + "/build/sync/" + start.AddDate(0, 0, i).Format("2006_01_02.log"))
		}

	case "monthly":
		start, end := lastPeriod(t, month)
		days := int(math.Abs(math.Round(start.Sub(end).Hours() / 24)))

		for i := 1; i <= days; i++ {
			deleteFile(BasePath() + "/build/sync/" + start.AddDate(0, 0, i).Format("2006_01_02.log"))
		}

	case "quarterly":
		start, end := lastPeriod(t, quarter)
		days := int(math.Abs(math.Round(start.Sub(end).Hours() / 24)))

		for i := 1; i <= days; i++ {
			deleteFile(BasePath() + "/build/sync/" + start.AddDate(0, 0, i).Format("2006_01_02.log"))
		}

	case "yearly":
		start, end := lastPeriod(t, fullYear)
		days := int(math.Abs(math.Round(start.Sub(end).Hours() / 24)))

		for i := 1; i <= days; i++ {
			deleteFile(BasePath() + "/build/sync/" + start.AddDate(0, 0, i).Format("2006_01_02.log"))
		}
	}
}
