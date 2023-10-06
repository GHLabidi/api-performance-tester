package models

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"github.com/joho/godotenv"
)

type BenchmarkData struct {
	TestStartTime        int64         `json:"test_start_time"`
	TestUniqueName       string        `json:"test_unique_name"`
	TestDisplayName      string        `json:"test_display_name"`
	TestDescription      string        `json:"test_description"`
	DataFolder           string        `json:"data_folder"`
	TestMode             string        `json:"test_mode"`
	RequestURL           string        `json:"request_url"`
	ConcurrentRequests   int           `json:"concurrent_requests"`
	SleepBetweenRequests float64       `yaml:"SleepBetweenRequests"`
	TestDuration         int           `json:"test_duration"`
	TotalRequests        int           `json:"total_requests"`
	SuccessfulRequests   int           `json:"successful_requests"`
	FailedRequests       int           `json:"failed_requests"`
	RequestsPerSecond    float64       `json:"requests_per_second"`
	QueryDurationStats   DurationStats `json:"query_duration_stats"`
	RequestDurationStats DurationStats `json:"request_duration_stats"`
}

var (
	rawData []StatData
)

// function to create a new BenchmarkData struct and calculate the stats
func NewBenchmarkData(test Test, stats []StatData, failedRequests int, testStartTime int64) BenchmarkData {
	rawData = stats
	// load environment variables
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	DataFolder := os.Getenv("DATA_FOLDER")

	var benchmarkData BenchmarkData
	benchmarkData.TestStartTime = testStartTime
	benchmarkData.TestUniqueName = test.TestUniqueName
	benchmarkData.TestDisplayName = test.TestDisplayName
	benchmarkData.TestDescription = test.TestDescription
	benchmarkData.DataFolder = DataFolder + "/" + test.TestUniqueName + "/"
	benchmarkData.TestMode = test.TestMode
	benchmarkData.RequestURL = test.RequestURL
	benchmarkData.ConcurrentRequests = test.ConcurrentRequests
	benchmarkData.SleepBetweenRequests = test.SleepBetweenRequests
	benchmarkData.TestDuration = test.TestDuration
	benchmarkData.TotalRequests = len(stats) + failedRequests
	benchmarkData.SuccessfulRequests = len(stats)
	benchmarkData.FailedRequests = failedRequests
	benchmarkData.RequestsPerSecond = float64(len(stats)+failedRequests) / float64(test.TestDuration)
	benchmarkData.QueryDurationStats = CalculateStatsFromColumn(stats, "query_duration")
	benchmarkData.RequestDurationStats = CalculateStatsFromColumn(stats, "request_duration")
	return benchmarkData
}

// save the benchmark data to a json file
func (b BenchmarkData) Save() error {
	// TODO check if data is valid
	// create the data folder if it doesn't exist
	if _, err := os.Stat(b.DataFolder); os.IsNotExist(err) {
		os.Mkdir(b.DataFolder, 0755)
	}
	fmt.Println("Saving benchmark data to", b.DataFolder+"benchmark.json")
	// create the file
	file, err := os.Create(b.DataFolder + "benchmark.json")
	if err != nil {
		return err
	}
	defer file.Close()

	// write the json to the file
	err = json.NewEncoder(file).Encode(b)
	if err != nil {
		return err
	}

	return nil
}

// save raw data to csv file
func (b BenchmarkData) SaveRawData() error {
	// TODO check if b is valid
	if len(rawData) == 0 {
		return fmt.Errorf("no raw data to save")
	}

	// create output folder
	if _, err := os.Stat(b.DataFolder); os.IsNotExist(err) {
		os.Mkdir(b.DataFolder, 0777)
	}

	// create csv file
	csvFile, err := os.Create(b.DataFolder + "data.csv")
	if err != nil {
		return err
	}
	defer csvFile.Close()

	// write CSV records
	for _, statDataItem := range rawData {
		_, err := csvFile.WriteString(fmt.Sprintf("%d,%d,%d,%d\n", statDataItem.StartTimestamp, statDataItem.EndTimestamp, statDataItem.QueryDuration, statDataItem.RequestDuration))
		if err != nil {
			return err
		}
	}
	csvFile.Sync()

	return nil
}

func (b BenchmarkData) GenerateReport() error {
	pythonScript := "python/generate_report.py"

	cmd := exec.Command("python", pythonScript, b.TestUniqueName)
	// get working directory
	cmd.Dir = os.Getenv("PWD")
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error running python script")
		fmt.Println(err)
		return err
	}
	return nil

}
