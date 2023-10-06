package handlers

import (
	"net/http"
	"os"
	"os/exec"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func GetComparisonHandler(w http.ResponseWriter, r *http.Request) {

	// load environment variables
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	// comparison folder is where comparison reports are saved
	comparisonFolder := os.Getenv("COMPARISON_FOLDER")

	// get folder names from URL
	vars := mux.Vars(r)
	folder1 := vars["folder1"]
	folder2 := vars["folder2"]

	// check if comparison file exists either for folder1 vs folder2 or folder2 vs folder1
	comparisonFile1 := folder1 + "_vs_" + folder2 + "_comparison_report.html"
	comparisonFile2 := folder2 + "_vs_" + folder1 + "_comparison_report.html"
	comparisonPath1 := comparisonFolder + "/" + comparisonFile1
	comparisonPath2 := comparisonFolder + "/" + comparisonFile2
	// check if comparison file exists
	if _, err := os.Stat(comparisonPath1); err == nil {
		// read comparison file and return it
		comparisonReport, err := os.ReadFile(comparisonPath1)
		if err != nil {
			http.Error(w, "Error reading comparison report", http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "text/html")

		w.Write(comparisonReport)
		return
	}
	if _, err := os.Stat(comparisonPath2); err == nil {
		// read comparison file and return it
		comparisonReport, err := os.ReadFile(comparisonPath2)
		if err != nil {
			http.Error(w, "Error reading comparison report", http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type", "text/html")
		w.Write(comparisonReport)

		return
	}

	// comparison file does not exist
	// create it by calling the compare.py script
	cmd := exec.Command("python", "python/compare.py", folder1, folder2)
	err = cmd.Run()
	if err != nil {
		http.Error(w, "Error running Python script", http.StatusInternalServerError)
		return
	}

	// return comparison file
	comparisonReport, err := os.ReadFile(comparisonPath1)
	if err != nil {
		http.Error(w, "Error reading comparison report", http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "text/html")
	w.Write(comparisonReport)

}
