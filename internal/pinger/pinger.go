package pinger

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

var client http.Client = http.Client{Timeout: 3 * time.Second}

func PingRegion(region Region, times int) Result {

	// Warmup to avoid DNS & TLS time in the next call
	if err := helper(region.Url); err != nil {
		return Result{Region: region, Latency: -1, Status: err.Error()}
	}
	// Calculate the latency
	var totalDuration time.Duration
	for i := 0; i < times; i++ {
		startTime := time.Now()
		if err := helper(region.Url); err != nil {
			return Result{Region: region, Latency: -1, Status: err.Error()}
		}
		totalDuration += time.Since(startTime)
	}

	avgDuration := totalDuration / time.Duration(times)

	return Result{Region: region, Latency: avgDuration.Milliseconds(), Status: "success"}

}

func helper(u string) error {
	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return err
	}
	res, err := client.Do(req)
	if err != nil {
		if err.(*url.Error).Timeout() {
			return errors.New("timeout")
		}
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP %d", res.StatusCode)
	}
	io.Copy(io.Discard, res.Body)
	return nil
}
