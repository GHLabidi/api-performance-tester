package handlers

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func GetBenchmarkByIdHandler(w http.ResponseWriter, r *http.Request) {
	// load environment variables
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	dataFolder := os.Getenv("DATA_FOLDER")

	if dataFolder == "" {
		panic("DATA_FOLDER environment variable not set")
	}

	// get id from url
	id := mux.Vars(r)["id"]

	// check if id is empty
	if id == "" {
		http.Error(w, "Test id not specified", http.StatusBadRequest)
		return
	}

	testDataFolder := fmt.Sprintf("%s/%s", dataFolder, id)
	// check if test data folder exists
	if _, err := os.Stat(testDataFolder); os.IsNotExist(err) {
		http.Error(w, "Test not found", http.StatusNotFound)
		return
	}

	// check if report.html file exists
	reportFile := fmt.Sprintf("%s/report.html", testDataFolder)
	if _, err := os.Stat(reportFile); os.IsNotExist(err) {
		http.Error(w, "Report not found", http.StatusNotFound)
		return
	}

	// read report.html file and write to response
	reportData, err := os.ReadFile(reportFile)
	if err != nil {
		http.Error(w, "Error reading report", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.Write(reportData)
}
