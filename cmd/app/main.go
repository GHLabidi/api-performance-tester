package main

import (
	"net/http"
	"os"

	"github.com/GHLabidi/api-performance-tester/internal/handlers"
	"github.com/GHLabidi/api-performance-tester/internal/httpbenchmark"
	"github.com/GHLabidi/api-performance-tester/internal/models"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	// load environment variables
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	testsFilePath := os.Getenv("TESTS_FILE_PATH")
	port := os.Getenv("PORT")

	// load tests from yaml file
	tests, err := models.LoadTestsFromFile(testsFilePath)
	if err != nil {
		panic(err)
	}
	// run benchmarks, generate reports and save results.
	httpbenchmark.RunTests(tests)
	// start http server
	startServer(port)

}

// TODO move this to a separate package
func startServer(port string) {
	// create router
	router := mux.NewRouter()
	// register handlers
	router.HandleFunc("/benchmarks", handlers.ListBenchmarksHandler).Methods("GET")                 // return list of completed benchmarks
	router.HandleFunc("/benchmarks/{id}", handlers.GetBenchmarkByIdHandler).Methods("GET")          // return benchmark report
	router.HandleFunc("/compare", handlers.CompareHandler).Methods("GET")                           // a form to compare two benchmarks
	router.HandleFunc("/compare", handlers.CompareHandler).Methods("POST")                          // redirect to comparison report
	router.HandleFunc("/compare/{folder1}/{folder2}", handlers.GetComparisonHandler).Methods("GET") // return comparison report

	err := http.ListenAndServe(":"+port, router)
	if err != nil {
		panic(err)
	}
}
