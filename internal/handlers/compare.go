package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func CompareHandler(w http.ResponseWriter, r *http.Request) {
	// check if request is POST
	if r.Method == "POST" {
		// get selected folders from form data
		folder1 := r.FormValue("folder1")
		folder2 := r.FormValue("folder2")

		// redirect to comparison file
		http.Redirect(w, r, fmt.Sprintf("/compare/%s/%s", folder1, folder2), http.StatusSeeOther)
		return
	}
	// request is GET

	// load environment variables
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	dataFolder := os.Getenv("DATA_FOLDER")
	comparisonFolder := os.Getenv("COMPARISON_FOLDER")

	if dataFolder == "" || comparisonFolder == "" {
		panic("DATA_FOLDER or COMPARISON_FOLDER environment variable not set")
	}

	// get list of folders in data folder
	folders, err := os.ReadDir(dataFolder)
	if err != nil {
		http.Error(w, "Error reading data folder", http.StatusInternalServerError)
		return
	}

	// generate HTML template with select options for folders
	var options []string
	for _, folder := range folders {
		if folder.IsDir() {
			options = append(options, folder.Name())
		}
	}
	// parse template
	tmpl, err := template.ParseFiles("internal/templates/compare.html")
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}

	data := struct {
		Options []string
	}{options}

	// execute template
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Error executing template", http.StatusInternalServerError)
		return
	}
}
