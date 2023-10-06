# API Performance Tester
## Description
API Performance Tester is a load testing tool for APIs written in Go and Python. It uses a YAML file to define test parameters, such as test name, description, test mode, test duration, and request URL etc... This makes it easy to configure and run tests for different APIs and/or with different configurations. It also provides a web user interface to view and compare the results of the tests.
## How to run
### Prerequisites
- Go 1.14 or higher (needed for the go modules)
- Redis Server
### Steps
1. Clone the repository
2. Place all your files in the data directory
3. Copy and edit the .env.example file to .env
4. Add your tests in config/tests.yaml
5. Two options to run the go server
   - Option 1:
        - Run `go run cmd/app/main.go` in the root directory
    - Option 2:
        1. Run `go build -o app` in the root directory
        2. Run `./app`
6. The program will run the tests and start a web server to view and compare the results.
## Adding your own tests
### Prerequisites
Your API endpoint must return a json object with a "QueryDuration" field (representing the time it took to execute the query in nanoseconds). The rest of the fields are up to you.
- **Sample Response**
```json
{
    "QueryDuration": 123456789,
    "Results": [
        {
            "id": 1,
            "name": "John Doe",
            "age": 25
        },
        {
            "id": 2,
            "name": "Jane Doe",
            "age": 23
        }
    ]
}
```

### Steps
1. Add your tests in config/tests.yaml
```yaml
- TestUniqueName: simple_search
  TestDisplayName: Simple Search
  TestDescription: Searching for the word "amazing" with simple search
  RequestURL: http://localhost:8080/search?word=amazing&searchMode=simple
  RequestType: GET
  TestMode: continious # continious or concurrent. If concurrent, you must specify the number of concurrent requests.
  # ConcurrentRequests: 100 # Only needed if TestMode is concurrent. If not specified, it will default to 10.
  TestDuration: 10
# Repeat for other tests
```
2. Run the program. It will run the tests and start a web server to view and compare the results.
## Web User Interface
### /benchmarks
Lists all the benchmarks that have been run. Clicking on a benchmark will take you to the benchmark page.
### /benchmarks/{benchmark_name}
Returns the report for the benchmark with the name {benchmark_name}.
### /compare
A page to compare two benchmarks. It will generate a comparison report and redirect you to the comparison page.
### /compare/{benchmark_name_1}/{benchmark_name_2}
Returns the comparison report for the two benchmarks with the names {benchmark_name_1} and {benchmark_name_2}.



