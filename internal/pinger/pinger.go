package pinger

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

var client http.Client = http.Client{Timeout: 3 * time.Second}

func PingRegion(region Region, times int64) Result {

	// Warmup to avoid DNS & TLS time in the next call
	if err := helper(region.Url); err != nil {
		return Result{Region: region, Latency: 0, Status: err.Error()}
	}
	// Calculate the latency
	var totalDuration time.Duration
	for i := 0; i < int(times); i++ {
		startTime := time.Now()
		if err := helper(region.Url); err != nil {
			return Result{Region: region, Latency: 0, Status: err.Error()}
		}
		totalDuration += time.Since(startTime)
	}

	avgDuration := totalDuration / time.Duration(times)

	return Result{Region: region, Latency: avgDuration.Milliseconds(), Status: "success"}

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
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP %d", res.StatusCode)
	}
	io.Copy(io.Discard, res.Body)
	return nil
}
