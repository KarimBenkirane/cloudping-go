package pinger

import (
	"encoding/json"
	"os"
	"slices"
)

type Regions []Region

func loadRegions() (Regions, error) {
	var regions Regions
	data, err := os.ReadFile("regions.json")
	if err != nil {
		return regions, err
	}
	err = json.Unmarshal(data, &regions)
	if err != nil {
		return regions, err
	}
	return regions, nil
}

func filterRegions(providers []string, codes []string) ([]Region, error) {
	allRegions, err := loadRegions()
	if err != nil {
		return nil, err
	}

	if len(providers) == 0 && len(codes) == 0 {
		return allRegions, nil
	}

	result := make([]Region, 0, len(allRegions))

	for _, region := range allRegions {
		pMatch := len(providers) == 0 || slices.Contains(providers, region.Provider)
		cMatch := len(codes) == 0 || slices.Contains(codes, region.Code)

		if pMatch && cMatch {
			result = append(result, region)
		}
	}

	return result, nil
}
