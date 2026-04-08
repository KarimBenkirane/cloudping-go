package pinger

import (
	_ "embed"
	"encoding/json"
	"slices"
	"strings"
	"sync"
)

//go:embed regions.json
var regionsJson string

type Region struct {
	Provider string
	Name     string
	Code     string
	Url      string
}

type Regions []Region

type Result struct {
	Region  Region
	Latency int64
	Status  string
}

type ResultStore struct {
	mu      sync.RWMutex
	results []Result
}

func (rs *ResultStore) Add(result Result) {
	rs.mu.Lock()
	defer rs.mu.Unlock()
	rs.results = append(rs.results, result)
}

func (rs *ResultStore) All() []Result {
	rs.mu.RLock()
	defer rs.mu.RUnlock()

	res := make([]Result, len(rs.results))
	copy(res, rs.results)
	return res
}

func loadRegions() (Regions, error) {
	dec := json.NewDecoder(strings.NewReader(regionsJson))
	var r Regions
	if err := dec.Decode(&r); err != nil {
		return nil, err
	}
	return r, nil
}

func FilterRegions(providers []string, codes []string) (Regions, error) {
	allRegions, err := loadRegions()
	if err != nil {
		return nil, err
	}

	if len(providers) == 0 && len(codes) == 0 {
		return allRegions, nil
	}

	result := make(Regions, 0, len(allRegions))

	for _, region := range allRegions {
		pMatch := len(providers) == 0 || slices.Contains(providers, region.Provider)
		cMatch := len(codes) == 0 || slices.Contains(codes, region.Code)

		if pMatch && cMatch {
			result = append(result, region)
		}
	}

	return result, nil
}
