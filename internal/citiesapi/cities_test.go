package citiesapi

import (
	"fmt"
	"log"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildCities(t *testing.T) {
	testCases := []struct {
		description   string
		filePath      string
		expectedErr   error
		expectedCount int
	}{
		{
			description:   "Success: Cities are built",
			filePath:      "data/mocks/mock_cities_good.json",
			expectedCount: 5,
			expectedErr:   nil,
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			t.Log(tc.description)
			t.Parallel()

			rootDir, err := getRootDir()
			if err != nil {
				log.Fatalf("Broken Test: %v", err)
			}

			absolutePath := filepath.Join(rootDir, tc.filePath)

			cities, err := BuildCities(absolutePath)
			if tc.expectedErr != nil {
				assert.Error(t, err, tc.description)
				return
			}

			assert.NoError(t, err, tc.description)
			assert.Equal(t, tc.expectedCount, len(cities))
		})
	}
}

// required to get the mocks file path while running the test inside of the citiesapi directory
func getRootDir() (string, error) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("failed to get caller information")
	}

	dir := filepath.Dir(filename)
	for i := 0; i < 10; i++ {
		if filepath.Base(dir) == "Iontra_TH" {
			break
		}

		dir = filepath.Dir(dir)
	}

	return dir, nil
}
