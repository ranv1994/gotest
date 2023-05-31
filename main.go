package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"

	_ "github.com/go-sql-driver/mysql"
)

// Define the JSON data you want to send
// var jsonData = []byte(`{"order_type": "value1", "id": 59}`)

// Define an array of endpoint URLs
var urls = []string{
	"http://localhost/user_dashboard/subscription/mt_test_webhook",
	"http://localhost/user_dashboard/subscription/mt_test_webhook",
	"http://localhost/user_dashboard/subscription/mt_test_webhook",
	"http://localhost/user_dashboard/subscription/mt_test_webhook",
	"http://localhost/user_dashboard/subscription/mt_test_webhook",
	"http://localhost/user_dashboard/subscription/mt_test_webhook",
	"http://localhost/user_dashboard/subscription/mt_test_webhook",
	"http://localhost/user_dashboard/subscription/mt_test_webhook",
	"http://localhost/user_dashboard/subscription/mt_test_webhook",
	"http://localhost/user_dashboard/subscription/mt_test_webhook",
	"http://localhost/user_dashboard/subscription/mt_test_webhook",
	"http://localhost/user_dashboard/subscription/mt_test_webhook",
	"http://localhost/user_dashboard/subscription/mt_test_webhook",
	"http://localhost/user_dashboard/subscription/mt_test_webhook",
	"http://localhost/user_dashboard/subscription/mt_test_webhook",
	"http://localhost/user_dashboard/subscription/mt_test_webhook",
	"http://localhost/user_dashboard/subscription/mt_test_webhook",
	"http://localhost/user_dashboard/subscription/mt_test_webhook",
	"http://localhost/user_dashboard/subscription/mt_test_webhook",
	"http://localhost/user_dashboard/subscription/mt_test_webhook",
	"http://localhost/user_dashboard/subscription/mt_test_webhook",
	"http://localhost/user_dashboard/subscription/mt_test_webhook",
	"http://localhost/user_dashboard/subscription/mt_test_webhook",
	"http://localhost/user_dashboard/subscription/mt_test_webhook",
	// ... (all 1000 URLs)
}

// Define an array to store the client endpoint URLs
var clientUrls []string

// Function to fetch client endpoint URLs from the database
func fetchClientUrlsFromDatabase() error {
	// Establish a database connection
	db, err := sql.Open("mysql", "username:password@tcp(localhost:3306)/database")
	if err != nil {
		return err
	}
	defer db.Close()

	// Execute the query to fetch the endpoint URLs
	rows, err := db.Query("SELECT url FROM endpoints")
	if err != nil {
		return err
	}
	defer rows.Close()

	// Iterate over the query results and populate the clientUrls array
	for rows.Next() {
		var url string
		err = rows.Scan(&url)
		if err != nil {
			return err
		}
		clientUrls = append(clientUrls, url)
	}

	err = rows.Err()
	if err != nil {
		return err
	}

	return nil
}

// Function to send JSON data to multiple endpoints concurrently
func sendJsonDataToMultipleEndpoints(data map[string]interface{}, wg *sync.WaitGroup) {
	defer wg.Done()

	// Convert data map to JSON bytes
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error marshaling JSON data: %v\n", err)
		return
	}
	for _, url := range urls {
		resp, err := http.Post(url, "application/json", strings.NewReader(string(jsonBytes)))
		if err != nil {
			log.Printf("Error sending data to %s: %v\n", url, err)
			continue
		}
		defer resp.Body.Close()

		log.Printf("Data sent to %s. Response: %d\n", url, resp.StatusCode)
	}
}

// Handler for triggering the data sending process
func triggerDataSending(w http.ResponseWriter, r *http.Request) {
	// Fetch JSON data from another API
	resp, err := http.Get("https://rest.entitysport.com/v2/matches/63029/live?token=dbee2220638adcb5a972ac42e3771c07")
	if err != nil {
		log.Fatalf("Error fetching JSON data: %v\n", err)
	}
	defer resp.Body.Close()

	// Read the JSON response into a byte slice
	jsonData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading JSON data: %v\n", err)
	}

	// Convert the JSON data to a map[string]interface{}
	var data map[string]interface{}
	err = json.Unmarshal(jsonData, &data)
	if err != nil {
		log.Fatalf("Error unmarshaling JSON data: %v\n", err)
	}

	// Fetch client endpoint URLs from the database
	err2 := fetchClientUrlsFromDatabase()
	if err2 != nil {
		log.Fatalf("Error fetching client URLs: %v\n", err2)
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go sendJsonDataToMultipleEndpoints(data, &wg)
	wg.Wait()

	fmt.Fprintf(w, "Data sending process completed")
}

func main() {
	http.HandleFunc("/trigger-data-sending", triggerDataSending)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
