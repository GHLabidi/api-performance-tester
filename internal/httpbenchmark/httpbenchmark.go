package httpbenchmark

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/GHLabidi/api-performance-tester/internal/models"
)

var (
	// TestStartTime is the time when the test started
	TestStartTime int64
)

func RunSingleTest(test models.Test) error {
	TestStartTime = time.Now().UnixNano()
	stats := []models.StatData{}
	failedRequests := 0

	// print test info
	fmt.Println("Test:", test.TestUniqueName)

	// check test mode
	switch test.TestMode {
	case "continious":
		fmt.Println("Starting continuous test.")
		// run continuous test
		stats, failedRequests = PerformContiniusTest(test)
		break
	case "concurrent":
		fmt.Println("Starting concurrent test.")
		// run concurrent test
		stats, failedRequests = PerformConcurrentTest(test)
		break
	default:
		fmt.Println("Test mode not specified. Defaulting to continuous")
		// run continuous test
		stats, failedRequests = PerformContiniusTest(test)
	}
	fmt.Println("Test finished.")
	fmt.Println("Total requests:", len(stats)+failedRequests)
	fmt.Println("Failed requests:", failedRequests)

	// create benchmark data, save it and generate report
	fmt.Println("Creating benchmark data.")
	benchmark := models.NewBenchmarkData(test, stats, failedRequests, TestStartTime)
	fmt.Println("Saving benchmark data.")
	benchmark.Save()
	fmt.Println("Saving raw data.")
	benchmark.SaveRawData()
	fmt.Println("Generating Report.")
	benchmark.GenerateReport()
	fmt.Println("Done. You can now view the results in: http://localhost:8081/benchmarks/" + test.TestUniqueName)

	return nil
}

// RunTests runs a list of tests
func RunTests(tests []models.Test) error {
	for _, test := range tests {
		RunSingleTest(test)
	}

	return nil
}

func PerformContiniusTest(test models.Test) ([]models.StatData, int) {

	var stats []models.StatData
	var failedRequests int = 0
	client := &http.Client{}
	// keep sending requests until test duration is reached
	for {
		// check if test duration is reached
		if time.Now().UnixNano()-TestStartTime > int64(test.TestDuration)*int64(time.Second) {
			break
		}
		// prepare the request
		req, err := http.NewRequest(test.RequestType, test.RequestURL, nil)
		if err != nil {
			fmt.Println("Error creating request.")
			fmt.Println(err)
			failedRequests++
			continue
		}
		req.Header.Set("Content-Type", "application/json")

		startTime := time.Now()
		// send request
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error sending request.")
			fmt.Println(err)
			failedRequests++
			continue
		}
		elapsed := time.Since(startTime)
		// extract duration from response
		fieldToExtract := "QueryDuration" // TODO make this configurable
		query_duration, err := extractDurationFromResponse(resp, fieldToExtract)
		if err != nil {
			fmt.Println("Error extracting duration.")
			fmt.Println(err)
			failedRequests++
			continue
		}
		// append stats
		stats = append(stats, models.StatData{
			StartTimestamp:  startTime.UnixNano(),
			EndTimestamp:    time.Now().UnixNano(),
			QueryDuration:   query_duration,
			RequestDuration: elapsed.Nanoseconds(),
		})

	}

	return stats, failedRequests

}

func PerformConcurrentTest(test models.Test) ([]models.StatData, int) {
	if test.ConcurrentRequests == 0 {
		fmt.Println("Concurrent requests number is not specified. Defaulting to 10.")
		test.ConcurrentRequests = 10
	}
	type ChannelResult struct {
		Stats          []models.StatData
		FailedRequests int
	}
	// create a buffered channel to receive results from goroutines
	ch := make(chan ChannelResult, test.ConcurrentRequests)
	// create a wait group to wait for all goroutines to finish
	wg := &sync.WaitGroup{}

	// start test.ConcurrentRequests goroutines
	for i := 0; i < test.ConcurrentRequests; i++ {
		wg.Add(1)
		// start a goroutine
		go func() {
			defer wg.Done()

			tmpStats := []models.StatData{}
			failedRequests := 0
			// keep sending requests until test duration is reached
			for {
				if time.Now().UnixNano()-TestStartTime > int64(test.TestDuration)*int64(time.Second) {
					break
				}
				stat, err := performSingleRequest(test)
				if err != nil {
					failedRequests++
					continue
				}
				tmpStats = append(tmpStats, stat)

			}
			// send results to channel
			ch <- ChannelResult{
				Stats:          tmpStats,
				FailedRequests: failedRequests,
			}

		}()

	}
	// wait for all goroutines to finish
	wg.Wait()

	// close the channel
	close(ch)

	// collect results from channel
	var stats []models.StatData
	var failedRequests int = 0
	for duration := range ch {
		stats = append(stats, duration.Stats...)
		failedRequests += duration.FailedRequests
	}
	return stats, failedRequests

}

func performSingleRequest(test models.Test) (models.StatData, error) {
	client := &http.Client{}
	req, err := http.NewRequest(test.RequestType, test.RequestURL, nil)
	if err != nil {
		fmt.Println("Error creating request.")
		fmt.Println(err)
		return models.StatData{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	startTime := time.Now()
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request.")
		fmt.Println(err)
		return models.StatData{}, err
	}
	elapsed := time.Since(startTime)
	fieldToExtract := "QueryDuration"
	query_duration, err := extractDurationFromResponse(resp, fieldToExtract)
	if err != nil {
		fmt.Println("Error extracting duration.")
		fmt.Println(err)
		return models.StatData{}, err
	}
	return models.StatData{
		StartTimestamp:  startTime.UnixNano(),
		EndTimestamp:    time.Now().UnixNano(),
		QueryDuration:   query_duration,
		RequestDuration: elapsed.Nanoseconds(),
	}, nil
}

func extractDurationFromResponse(resp *http.Response, fieldToExtract string) (int64, error) {
	var query_duration int64
	// get the response body
	var responseData map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&responseData); err != nil {
		fmt.Println("Error decoding response.")
		fmt.Println(err)
		return 0, err
	}
	// check if the field exists
	if responseData[fieldToExtract] == nil {
		fmt.Println("Error extracting field.")
		return 0, fmt.Errorf("Error extracting field %s", fieldToExtract)
	}
	// convert the value to int64 and check if it's valid
	if value, ok := responseData[fieldToExtract].(float64); ok {
		query_duration = int64(value)
	} else {
		fmt.Println("Error converting value.")
		return 0, fmt.Errorf("Error converting value %v", responseData[fieldToExtract])
	}

	return query_duration, nil
}
