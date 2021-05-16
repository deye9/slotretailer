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
	"reflect"
	"strconv"
	"strings"
	"syscall"
	"time"
)

// APIlinks is a collection of all needed API links
var APIlinks = make(map[string]string)
var apikeys = [...]string{"customers", "orders", "transfers", "transferRequests"}

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

	if LocalStore.CreditCardAPI == "" || LocalStore.ProductsAPI == "" || LocalStore.CustomersAPI == "" || LocalStore.SapKey == "" {
		CheckError("LocalStore endpoint missing.", errors.New("Missing endpoint from application"), false)
		return
	}

	APIlinks["orders"] = LocalStore.OrdersAPI                                                             // WIP
	APIlinks["stores"] = LocalStore.WarehousesAPI                                                         // Ready
	APIlinks["cheques"] = LocalStore.ChequesAPI                                                           // Ready
	APIlinks["products"] = LocalStore.ProductsAPI                                                         // Ready
	APIlinks["customers"] = LocalStore.CustomersAPI                                                       // POST ready and GET ready. How to get customers for other stores / across board.
	APIlinks["pricelist"] = LocalStore.PricelistAPI                                                       // Ready
	APIlinks["creditcards"] = LocalStore.CreditCardAPI                                                    // Ready
	APIlinks["cashaccounts"] = LocalStore.CashAccountAPI                                                  // Ready
	APIlinks["banktransfer"] = LocalStore.BankTransferAPI                                                 // Ready
	APIlinks["transferRequests"] = LocalStore.TransfersAPI                                                // Ready
	APIlinks["transfers"] = strings.Replace(LocalStore.TransfersAPI, "TransferRequests", "Transfers", -1) // Ready

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
	updateStatus()
	getTransfers()

	for key, link := range APIlinks {
		// Write the sync start details to the File System via a Goroutine.
		go WriteFile(BasePath()+"/build/sync/"+str+".log", []byte("Sync for "+key+" started at "+t.String()+"\n"))

		// Append the StoreID to the link
		if key == "products" {
			link += "?storeId=" + LocalStore.SapKey + "&pricelist=" + LocalStore.ProductPriceList
		} else if key == "transfers" {
			link += "?sourceStore=" + LocalStore.SapKey
		} else {
			link += "?storeId=" + LocalStore.SapKey
		}

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

		if url == "" {
			continue
		}

		// Append the StoreID to the link
		if value != "transfers" && value != "transferRequests" {
			url += "?storeId=" + LocalStore.SapKey
		}

		switch value {
		case "customers":
			SQLquery = "select id, cardname, address, phone, phone1, city, email, synced, created_by from customers where synced = false and deleted_at is null;"

		case "orders":
			SQLquery = "call getOrders()"

		case "transfers":
			SQLquery = `select id, fromwhs, towhs, comment, canceled, synced, status, created_by, DATE_FORMAT(created_at, '%Y-%m-%d') docdate, docentry, docnum, requestId,
			(select JSON_ARRAYAGG(JSON_OBJECT('itemcode', itemcode, 'itemname', itemname, 'quantity', quantity, 'serialnumber', serialnumber)) from transfereditems 
			where transferid = t.id) items from transfers t where deleted_at is null and status = 'Accepted' and synced = false;`

		case "transferRequests":
			SQLquery = `select id, fromwhs, towhs, comment, canceled, synced, status, created_by, DATE_FORMAT(created_at, '%Y-%m-%d') docdate, docentry, docnum, requestId,
			(select JSON_ARRAYAGG(JSON_OBJECT('itemcode', itemcode, 'itemname', itemname, 'quantity', quantity, 'serialnumber', serialnumber)) from transfereditems 
			where transferid = t.id) items from transfers t where deleted_at is null and status in ('pending', 'rejected') and synced = false;`
		}

		var rows *sql.Rows
		var columns []string

		if rows, err = Get(SQLquery); err != nil {
			CheckError("Error Getting "+value+" data for Endpoint: ", err, false)
			continue
		}

		if columns, err = rows.Columns(); err != nil {
			CheckError("Error getting Row columns from "+value+"Query.", err, false)
			continue
		}

		defer rows.Close()
		ConvertToJSON(rows, columns, url, key)
	}
	return nil
}

func getTransfers() (err error) {
	data := []byte{}
	var transferLinks = make(map[string]string)
	transferLinks["unprocessedDestination"] = LocalStore.TransfersAPI + "/Unprocessed/destination?destinationStore=" + LocalStore.SapKey

	for key, link := range transferLinks {
		var response interface{} = []Transfers{}
		data, err = httpget(link)

		if string(data) == "[]" {
			continue
		}

		if err = json.Unmarshal(data, &response); err != nil {
			CheckError("Error unmarshalling data for ["+key+"].", err, false)
			continue
		}

		cmd := structToInsert(response, "transfers")
		if strings.ToLower((key)) == "unprocesseddestination" {
			cmd = strings.ReplaceAll(strings.ToLower(cmd), "id", "requestid")
			cmd = strings.ReplaceAll(strings.ToLower(cmd), "pending", "Incoming")
			cmd = strings.ReplaceAll(strings.ToLower(cmd), "docdate", "created_at")
			cmd = strings.ReplaceAll(strings.ToLower(cmd), "transferrequestid", "transferid")
		}

		if err = Modify(cmd); err != nil {
			CheckError("Error saving HTTPGET for "+key+" result.", err, false)
			continue
		}
	}
	return
}

func updateStatus() (err error) {
	var rows *sql.Rows
	if rows, err = Get("select docentry, status  from transfers where lower(status) like 'finalize_%' and deleted_at is null and synced = false order by created_at desc;"); err != nil {
		CheckError("Error getting Transfer Request Finalized data for sync.", err, false)
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var docentry int
		var status string
		if err = rows.Scan(&docentry, &status); err != nil {
			CheckError("Error Scanning Transfer Request Details.", err, false)
			continue
		}

		url := APIlinks["transfers"] + "/UpdateStatus/" + strconv.Itoa(docentry) + "?status=" + status
		method := "PUT"

		client := &http.Client{}
		req, err := http.NewRequest(method, url, nil)
		if err != nil {
			CheckError("Error Parsing PUT data to URL: "+url, err, false)
			continue
		}

		res, err := client.Do(req)
		if err != nil {
			CheckError("Error PATCHING data to URL: "+url, err, false)
			continue
		}
		defer res.Body.Close()

		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			CheckError("Error from Endpoint "+url+" for PATCH payload sent. ", err, false)
			continue
		}

		if status == "200 OK" {
			Modify("UPDATE transfers set synced = true; ")
		} else {
			CheckError("Error from Endpoint "+url+" for Payload sent. ", errors.New("status is "+status+". "+string(data)), false)
		}
	}

	return
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
				resultMap[columns[i]] = fmt.Sprintf("%s", val)
			} else if strings.ToLower(columns[i]) == "synced" {
				resultMap[columns[i]] = true
			} else if strings.ToLower(columns[i]) == "returned" {
				switch fmt.Sprintf("%s", val) {
				case "0":
					resultMap[columns[i]] = false
				case "1":
					resultMap[columns[i]] = true
				}
			} else if strings.ToLower(columns[i]) == "canceled" {
				switch fmt.Sprintf("%s", val) {
				case "0":
					resultMap[columns[i]] = false
				case "1":
					resultMap[columns[i]] = true
				}
			} else if strings.ToLower(columns[i]) == "items" {
				resultMap[columns[i]] = interfaceToMap(val, "ordered items")
			} else if strings.ToLower(columns[i]) == "payments" {
				resultMap[columns[i]] = interfaceToMap(val, "payment")
			} else if strings.ToLower(columns[i]) == "comment" {
				if val == nil {
					resultMap[columns[i]] = ""
				} else {
					resultMap[columns[i]] = fmt.Sprintf("%s", val)
				}
			} else {
				resultMap[columns[i]] = fmt.Sprintf("%s", val)
			}
		}

		// for each database row / record, a map with the column names and row values is added to the allMaps slice
		allMaps = append(allMaps, resultMap)
	}

	// Marshal the map into a JSON string.
	payload, err := json.Marshal(allMaps)
	if err != nil {
		CheckError("Error Marshalling the map into a JSON string.", err, false)
		return
	}

	jsonStr := string(payload)
	if jsonStr != "null" {
		// Remove the last ", " from the ID string and generate the update command
		cmd := "UPDATE " + strings.Replace(key, "transferRequests", "transfers", -1) + " SET synced = true WHERE id IN (" + strings.TrimRight(id, ", ") + ");"
		_, _, _ = httppost(url, jsonStr, cmd)
	}
	return
}

func interfaceToMap(val interface{}, message string) (mapped []map[string]interface{}) {
	err := json.Unmarshal([]byte(fmt.Sprintf("%s", val)), &mapped)
	if err != nil {
		fmt.Println("Error unmarshalling ", message, " details for sync: ", err)
	}
	return
}

// httppost to post the data to the server
func httppost(url, payload, successcommand string) (status string, data []byte, err error) {
	method := "POST"
	requestBody := strings.NewReader(payload)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, requestBody)
	if err != nil {
		CheckError("Error Parsing data to URL: "+url, err, false)
		return
	}

	req.Header.Add("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		CheckError("Error POSTING data to URL: "+url, err, false)
		if strings.Contains(err.Error(), "wsarecv") {
			return
		}
	}

	defer res.Body.Close()

	status = res.Status
	data, err = ioutil.ReadAll(res.Body)
	if err != nil {
		CheckError("Error from Endpoint "+url+" for Payload sent. ", err, false)
		return
	}

	if status == "200 OK" {
		Modify(successcommand)
	} else {
		CheckError("Error from Endpoint "+url+" for Payload sent. ", errors.New("status is "+status+". "+string(data)), false)
	}

	fmt.Println("URL IS: ", url)
	fmt.Println("data IS: ", string(data))
	fmt.Println("status IS: ", status)
	return
}

func getAllData(key, link, str string) error {
	cmd := ""
	data, err := httpget(link)

	if string(data) == "" {
		return nil
	}

	var response interface{}

	if strings.ToLower(key) == "orders" {
		response = []Orders{}
	} else if strings.ToLower(key) == "cheques" {
		response = []Cheques{}
	} else if strings.ToLower(key) == "banktransfer" {
		response = []BankTransfer{}
	} else if strings.ToLower(key) == "products" {
		response = []Products{}
	} else if strings.ToLower(key) == "customers" {
		response = []Customers{}
	} else if strings.ToLower(key) == "pricelist" {
		response = []PriceList{}
	} else if strings.ToLower(key) == "creditcards" {
		response = []CreditCards{}
	} else if strings.ToLower(key) == "cashaccounts" {
		response = []CashAccounts{}
	} else if strings.ToLower(key) == "transfers" {
		response = []Transfers{}
	}

	if err = json.Unmarshal(data, &response); err != nil {
		return err
	}

	if reflect.TypeOf(response).Kind().String() == "map" {
		return nil
	}

	if len(response.([]interface{})) <= 0 {
		// Write the sync start details to the File System via a Goroutine.
		go WriteFile(BasePath()+"/build/sync/"+str+".log", []byte("Sync for "+key+" finished at "+time.Now().String()+"\n"))
		return nil
	}

	switch strings.ToLower(key) {
	case "orders":
	case "customers":
	case "transfers":
	default:
		cmd = "truncate table " + key + "; "
	}

	if strings.ToLower(key) == "transfers" {
		cmd += structToUpdate(response, key)
	} else {
		cmd += structToInsertUpdate(response, key)
	}

	if err = Modify(cmd); err != nil {
		CheckError("Error saving HTTPGET for "+key+" result.", err, false)
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
