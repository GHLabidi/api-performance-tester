package models

type StatData struct {
	StartTimestamp  int64 `json:"start_timestamp"`
	EndTimestamp    int64 `json:"end_timestamp"`
	QueryDuration   int64 `json:"query_duration"`
	RequestDuration int64 `json:"request_duration"`
}
