package models

import "time"

type Reply struct {
	Word          string
	Count         int
	Files         []string
	QueryDuration time.Duration
	SearchMode    string
}
