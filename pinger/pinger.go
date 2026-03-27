package pinger

import (
	"io"
	"net/http"
	"time"
)

var client http.Client = http.Client{Timeout: 3 * time.Second}

func pingRegion(region Region) Result {

	// Warmup to avoid DNS & TLS time in the next call
	err := helper(region.Url)
	if err != nil {
		return Result{Region: region, Latency: 0, Status: err.Error()}
	}

	// Calculate the latency
	startTime := time.Now()
	err = helper(region.Url)
	if err != nil {
		return Result{Region: region, Latency: 0, Status: err.Error()}
	}
	endTime := time.Since(startTime)

	return Result{Region: region, Latency: endTime.Milliseconds(), Status: "success"}

}

func helper(url string) error {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	io.Copy(io.Discard, res.Body) // Read the body until EOF to "re-use a persistent TCP connection to the server for a subsequent "keep-alive" request."
	return nil
}

func Ping(regions Regions) []Result {
	results := make([]Result, 0, len(regions))
	for _, region := range regions {
		result := pingRegion(region)
		results = append(results, result)
	}
	return results
}
