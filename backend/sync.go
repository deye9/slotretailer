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

	// APIlinks["orders"] = LocalStore.OrdersAPI
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
	fmt.Println("work started at ", t)
	getAllData()

	time.Sleep(2 * time.Second)
	fmt.Println("work finished at ", time.Now())
}

// check for connectivity. If err == `Get "https://5f63ea94363f0000162d9307.mockapi.io/api/v1/Products": dial tcp: lookup 5f63ea94363f0000162d9307.mockapi.io: no such host` then we have no internet

func getAllData() error {
	for key, link := range APIlinks {
		cmd := ""
		data, err := httpget(link)
		var response interface{}

		if strings.ToLower(key) == "orders" {
			response = []Orders{}
		} else if strings.ToLower(key) == "products" {
			response = []Products{}
		} else if strings.ToLower(key) == "customers" {
			response = []Customers{}
		}

		if err = json.Unmarshal(data, &response); err != nil {
			CheckError("Error unmarshalling data.", err, false)
			return err
		}
		cmd = structToReplace(response, key)

		if err = Modify(cmd); err != nil {
			CheckError("Error saving HTTPGET result. ", err, false)
			return err
		}
	}
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
