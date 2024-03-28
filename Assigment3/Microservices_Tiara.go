package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"sync"
	"time"
)

// Status represents the structure of the status JSON
type Status struct {
	Water int `json:"water"`
	Wind  int `json:"wind"`
}

// StatusInfo represents additional information about the status
type StatusInfo struct {
	StatusWater string `json:"statusWater"`
	StatusWind  string `json:"statusWind"`
}

var (
	statusLock sync.Mutex
	status     Status
)

func main() {
	// Start updating status every 15 seconds
	go updateStatus()

	// Serve status page
	http.HandleFunc("/", statusHandler)

	// Start HTTP server
	fmt.Println("Server is running on port 8080...")
	http.ListenAndServe(":8080", nil)
}

// updateStatus updates the status JSON every 15 seconds with random values
func updateStatus() {
	for {
		// Generate random values for water and wind
		statusLock.Lock()
		status.Water = rand.Intn(100) + 1
		status.Wind = rand.Intn(100) + 1
		statusLock.Unlock()

		// Write status to JSON file
		writeStatusToFile()

		// Wait for 15 seconds
		time.Sleep(15 * time.Second)
	}
}

// writeStatusToFile writes the current status to JSON file
func writeStatusToFile() {
	statusLock.Lock()
	defer statusLock.Unlock()

	data, err := json.MarshalIndent(status, "", "    ")
	if err != nil {
		fmt.Println("Error marshaling status:", err)
		return
	}

	err = ioutil.WriteFile("status.json", data, 0644)
	if err != nil {
		fmt.Println("Error writing status to file:", err)
		return
	}
}

// getStatusInfo returns additional information about the status
func getStatusInfo() StatusInfo {
	statusLock.Lock()
	defer statusLock.Unlock()

	var info StatusInfo

	// Determine status for water
	switch {
	case status.Water < 5:
		info.StatusWater = "Aman"
	case status.Water >= 5 && status.Water <= 8:
		info.StatusWater = "Siaga"
	default:
		info.StatusWater = "Bahaya"
	}

	// Determine status for wind
	switch {
	case status.Wind < 6:
		info.StatusWind = "Aman"
	case status.Wind >= 6 && status.Wind <= 15:
		info.StatusWind = "Siaga"
	default:
		info.StatusWind = "Bahaya"
	}

	return info
}

// statusHandler handles requests to the status page
func statusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	// Read status from file
	statusLock.Lock()
	defer statusLock.Unlock()

	file, err := os.Open("status.json")
	if err != nil {
		http.Error(w, "Failed to read status.", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		http.Error(w, "Failed to read status.", http.StatusInternalServerError)
		return
	}

	var currentStatus Status
	err = json.Unmarshal(data, &currentStatus)
	if err != nil {
		http.Error(w, "Failed to read status.", http.StatusInternalServerError)
		return
	}

	// Get additional status information
	info := getStatusInfo()

	// Generate HTML response
	html := fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<head>
			<meta http-equiv="refresh" content="15">
			<title>Status</title>
		</head>
		<body>
			<h1>Current Status</h1>
			<p>Water: %d meter</p>
			<p>Status Water: %s</p>
			<p>Wind: %d meter/s</p>
			<p>Status Wind: %s</p>
		</body>
		</html>
	`, currentStatus.Water, info.StatusWater, currentStatus.Wind, info.StatusWind)

	// Write response
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}
