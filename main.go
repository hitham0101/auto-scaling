package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const prometheusURL = "http://195.201.239.63:9090/api/v1/query"

// PrometheusResponse defines the structure of Prometheus query response
type PrometheusResponse struct {
	Status string `json:"status"`
	Data   struct {
		ResultType string `json:"resultType"`
		Result     []struct {
			Metric map[string]string `json:"metric"`
			Value  []interface{}     `json:"value"`
		} `json:"result"`
	} `json:"data"`
}

func getCPUUsage() (string, error) {
	// Query for CPU usage
	query := `100 - (avg(rate(node_cpu_seconds_total{mode="idle"}[1m])) * 100)`
	resp, err := http.Get(fmt.Sprintf("%s?query=%s", prometheusURL, query))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Read and parse the response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result PrometheusResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	// Extract CPU usage value
	if len(result.Data.Result) > 0 {
		value := result.Data.Result[0].Value[1]
		return fmt.Sprintf("CPU Usage: %s%%", value), nil
	}
	return "No data available", nil
}

func cpuUsageHandler(w http.ResponseWriter, r *http.Request) {
	cpuUsage, err := getCPUUsage()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching CPU usage: %v", err), http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(w, cpuUsage)
}

func main() {
	http.HandleFunc("/cpu-usage", cpuUsageHandler)
	log.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
