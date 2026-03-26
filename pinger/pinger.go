package pinger

import (
	"net/http"
	"time"
)

func getCurrentTime() time.Time {
	return time.Now()
}

func PingRegion(region Region) Result {
	startTime := getCurrentTime()
	_, err := http.Head(region.Url)
	endTime := time.Since(startTime).Milliseconds()
	if err != nil {
		return Result{Region: region, Latency: endTime, Status: err.Error()}
	}
	return Result{Region: region, Latency: endTime, Status: "success"}

}

func Ping(regions []Region) []Result {
	results := make([]Result, 0, len(regions))
	for _, region := range regions {
		result := PingRegion(region)
		results = append(results, result)
	}
	return results
}
