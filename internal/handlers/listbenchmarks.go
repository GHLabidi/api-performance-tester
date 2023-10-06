package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func ListBenchmarksHandler(w http.ResponseWriter, r *http.Request) {
	// load environment variables
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	dataFolder := os.Getenv("DATA_FOLDER")

	if dataFolder == "" {
		panic("DATA_FOLDER environment variable not set")
	}

	// get list of folders in data folder
	folders, err := os.ReadDir(dataFolder)
	if err != nil {
		http.Error(w, "Error reading data folder", http.StatusInternalServerError)
		return
	}
	fmt.Println("Folders:")
	fmt.Println(folders)

	// generate HTML template with links to benchmarks
	var benchmarkNames []string
	for _, folder := range folders {
		if folder.IsDir() {
			// check that folder contains report.html
			_, err := os.Stat(dataFolder + "/" + folder.Name() + "/report.html")
			if err != nil {
				continue
			}
			// valid benchmark folder, add to list
			benchmarkNames = append(benchmarkNames, folder.Name())
		}
	}

	tmpl, err := template.ParseFiles("internal/templates/listbenchmarks.html")
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, benchmarkNames)
	if err != nil {
		http.Error(w, "Error executing template", http.StatusInternalServerError)
		return
	}
}
