package service

import (
	"database/sql"
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
var apikeys = [...]string{"customers"} //, "orders", "transfers"}

// Sync will setup the cadence for sync btw the store and the server.
func Sync() {

	// rotateLogs if folder exists else create the sync folder
	if _, err := os.Stat(BasePath() + "/build/sync"); os.IsNotExist(err) {
		if err := os.Mkdir(BasePath()+"/build/sync", 0755); err != nil {
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
	APIlinks["banks"] = LocalStore.BanksAPI       // Ready
	APIlinks["stores"] = LocalStore.WarehousesAPI // Ready
	APIlinks["products"] = LocalStore.ProductsAPI	// Ready
	APIlinks["customers"] = LocalStore.CustomersAPI // POST ready and GET ready. How to get customers for other stores / across board.
	// APIlinks["transfers"] = LocalStore.TransfersAPI // Get for other products on a need to basis.

	duration := LocalStore.SyncInterval
	if duration == 0 {
		duration = 10
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

	sendData()

	for key, link := range APIlinks {
		// Write the sync start details to the File System via a Goroutine.
		go WriteFile(BasePath()+"/build/sync/"+str+".log", []byte("Sync for "+key+" started at "+t.String()+"\n"))

		// Append the StoreID to the link
		link += "?storeID=" + LocalStore.SapKey
		getAllData(key, link, str)
		time.Sleep(2 * time.Second)
	}
}

// sendData for Customers, Orders and Inventory Transfers
func sendData() (err error) {
	for _, value := range apikeys {
		SQLquery := ""
		key := value
		url := APIlinks[value]

		switch value {
		case "customers":
			SQLquery = "select id, cardname, address, phone, phone1, city, email, synced from customers where synced = false;"

		case "orders":
			SQLquery = "select id, cardname, address, phone, phone1, city, email from customers where synced = false;"

		case "transfers":
			SQLquery = "select id, cardname, address, phone, phone1, city, email from customers where synced = false;"
		}

		var rows *sql.Rows
		var columns []string

		if rows, err = Get(SQLquery); err != nil {
			CheckError("Error Getting customer data for Endpoint: ", errors.New(SQLquery), false)
			return
		}

		if columns, err = rows.Columns(); err != nil {
			CheckError("Error getting Row columns from Report.", err, false)
			return
		}

		defer rows.Close()
		ConvertToJSON(rows, columns, url, key)
	}
	return nil
}

// ConvertToJSON converts the database response to JSON.
func ConvertToJSON(rows *sql.Rows, columns []string, url, key string) (err error) {
	id := ""
	var allMaps []map[string]interface{}
	for rows.Next() {
		// Dynamic Result rows scanning.
		values := make([]interface{}, len(columns))
		pointers := make([]interface{}, len(columns))

		for i := range values {
			pointers[i] = &values[i]
		}
		err = rows.Scan(pointers...)

		resultMap := make(map[string]interface{})
		for i, val := range values {
			if strings.ToLower(columns[i]) == "id" {
				id += fmt.Sprintf("%s, ", val)
			}

			if strings.ToLower(columns[i]) == "synced" {
				resultMap[columns[i]] = true
				// switch fmt.Sprintf("%s", val) {
				// case "0":
				// 	resultMap[columns[i]] = false
				// case "1":
				// 	resultMap[columns[i]] = true
				// }
			} else {
				resultMap[columns[i]] = fmt.Sprintf("%s", val)
			}
		}

		// for each database row / record, a map with the column names and row values is added to the allMaps slice
		allMaps = append(allMaps, resultMap)
	}

	// Marshal the map into a JSON string.
	empData, err := json.Marshal(allMaps)
	if err != nil {
		CheckError("Error Marshalling the map into a JSON string.", err, false)
		return
	}

	jsonStr := string(empData)

	// Remove the last ", " from the ID string and generate the update command
	cmd := "UPDATE " + key + " SET synced = true WHERE id IN (" + strings.TrimRight(id, ", ") + ");"
	_, _, _ = httppost(url, jsonStr, cmd)

	// status, data, err := httppost(url, jsonStr, cmd)
	// fmt.Println("Status is: ", status)
	// fmt.Println("Response is: ", string(data))
	// fmt.Println("Error is: ", err)
	return
}

// httppost to post the data to the server
func httppost(url, payload, successcommand string) (status string, data []byte, err error) {
	method := "POST"
	requestBody := strings.NewReader(payload)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, requestBody)
	if err != nil {
		CheckError("Error POSTING data to URL: "+url, err, false)
	}

	req.Header.Add("Content-Type", "application/json")
	res, err := client.Do(req)
	defer res.Body.Close()

	status = res.Status
	data, err = ioutil.ReadAll(res.Body)

	if res.Status == "200 OK" {
		Modify(successcommand)
	}

	return
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
	} else if strings.ToLower(key) == "transfers" {
		response = []Transfers{}
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

// Data to send

// 1. Unsynced Customers
// 2. Unsynced Orders
// 3. Unsynced Transfers

// select * from customers where deleted_at is null and synced = false;
// select * from orders o  inner join ordereditems i on o.id = i.orderid where deleted_at is null and synced = false;
// select * from transfers t  inner join transfereditems i on t.id = i.transferid where deleted_at is null and synced = false;
