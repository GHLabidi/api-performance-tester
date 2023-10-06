package models

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Test struct {
	TestUniqueName       string  `yaml:"TestUniqueName"`
	TestDisplayName      string  `yaml:"TestDisplayName"`
	TestDescription      string  `yaml:"TestDescription"`
	RequestURL           string  `yaml:"RequestURL"`
	RequestType          string  `yaml:"RequestType"`
	RequestHeaders       string  `yaml:"RequestHeaders"`
	RequestBody          string  `yaml:"RequestBody"`
	TestMode             string  `yaml:"TestMode"`
	ConcurrentRequests   int     `yaml:"ConcurrentRequests"`
	SleepBetweenRequests float64 `yaml:"SleepBetweenRequests"`
	TestDuration         int     `yaml:"TestDuration"`
}

func LoadTestsFromFile(filepath string) ([]Test, error) {
	var tests []Test
	// print working directory
	wd, err := os.Getwd()
	if err != nil {
		return tests, err
	}
	println("Working directory:", wd)

	file, err := os.Open(filepath)
	if err != nil {
		return tests, err
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&tests)
	if err != nil {
		return tests, err
	}
	return tests, nil
}
