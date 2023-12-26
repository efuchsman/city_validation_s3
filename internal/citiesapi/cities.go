package citiesapi

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/efuchsman/city_validation_s3/internal/models"
	"github.com/pkg/errors"
)

// Slices for all of the files with bad json
type UnprocessableFiles struct {
	Files []string `json:"unprocessable_files"`
}

type CitiesMap struct {
	Cities map[string]*models.City
}

// Helper which builds a cities slice
func BuildCities(filePath string) ([]*models.City, error) {
	// Read the JSON file
	jsonData, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	// Create a slice to hold the cities
	var cities []*models.City

	// Unmarshal the JSON data into the cities slice
	err = json.Unmarshal(jsonData, &cities)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return cities, nil
}

// Builds TmpCity objects slice and finds unprocessable files for the tmp package
func BuildTMPCities() ([]*models.TmpCity, *UnprocessableFiles, error) {
	dirPath := "data/tmp"

	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return nil, nil, err
	}

	var cities []*models.TmpCity
	var badJson []string
	for _, file := range files {
		filePath := filepath.Join(dirPath, file.Name())

		jsonData, err := ioutil.ReadFile(filePath)
		if err != nil {
			fmt.Println("Error reading file:", err)
			continue
		}

		var jsonCities []*models.TmpCity
		err = json.Unmarshal(jsonData, &jsonCities)
		if err != nil {
			badJson = append(badJson, file.Name())
			continue
		}
		cities = append(cities, jsonCities...)
	}
	unprocessableFiles := &UnprocessableFiles{
		Files: badJson,
	}

	return cities, unprocessableFiles, nil
}

// Creates a Map from cities.json which is then stored in the client
func BuildCitiesMapFromCities() (*CitiesMap, error) {
	filePath := "data/cities.json"
	cities, err := BuildCities(filePath)
	if err != nil {
		return nil, err
	}

	citiesMap := &CitiesMap{
		Cities: make(map[string]*models.City),
	}
	for _, city := range cities {
		citiesMap.Cities[city.City] = city
	}

	return citiesMap, nil
}

func (c *client) ReturnCitiesMap() *CitiesMap {
	return c.cities
}

// Checks tmp city elements against the cities map
func (c *client) returnTMPCities() ([]*models.TmpCity, error) {
	var cities []*models.TmpCity
	for _, city := range c.tmpCities {

		key := &models.City{
			Latitude:     city.Latitude,
			Longitude:    city.Longitude,
			Geo:          city.Geo,
			City:         city.City,
			ProvinceIcon: city.ProvinceIcon,
			Province:     city.Province,
			CountryIcon:  city.CountryIcon,
			Country:      city.Country,
		}
		if c.cities.Cities[key.City].Latitude == key.Latitude &&
			c.cities.Cities[key.City].Longitude == key.Longitude &&
			c.cities.Cities[key.City].Geo == key.Geo &&
			c.cities.Cities[key.City].City == key.City &&
			c.cities.Cities[key.City].ProvinceIcon == key.ProvinceIcon &&
			c.cities.Cities[key.City].Province == key.Province &&
			c.cities.Cities[key.City].CountryIcon == key.CountryIcon &&
			c.cities.Cities[key.City].Country == key.Country {
			city.IsValid = true
		} else {
			city.IsValid = false
		}
		cities = append(cities, city)
	}

	return cities, nil
}

// Creates a slice for validated cities
func (c *client) ReturnValidTmpElements() ([]*models.TmpCity, error) {
	elements, err := c.returnTMPCities()
	if err != nil {
		return nil, err
	}

	validElements := make([]*models.TmpCity, 0)
	for _, e := range elements {
		if !e.IsValid {
			continue
		}
		validElements = append(validElements, e)
	}

	return validElements, nil
}

// Creates a slice for invalid cities
func (c *client) ReturnInvalidTmpElements() ([]*models.TmpCity, error) {
	elements, err := c.returnTMPCities()
	if err != nil {
		return nil, err
	}

	invalidElements := make([]*models.TmpCity, 0)
	for _, e := range elements {
		if e.IsValid {
			continue
		}
		invalidElements = append(invalidElements, e)
	}

	return invalidElements, nil
}

// Client call with cached unprocessable files
func (c *client) ReturnUnprocessableFiles() *UnprocessableFiles {
	return c.badFiles
}

// Builds a JSON output for a valid tmp cities input
func (c *client) CreateValidElementsJSON(elements []*models.TmpCity, fileName string) error {
	filePath := injectPath(fileName)
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	jsonData, err := json.MarshalIndent(elements, "", "    ")
	if err != nil {
		return err
	}

	_, err = file.Write(jsonData)
	if err != nil {
		return err
	}

	return nil
}

// Builds a JSON output for a invalid tmp cities input
func (c *client) CreateInvalidElementsJSON(elements []*models.TmpCity, fileName string) error {
	filePath := injectPath(fileName)
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	jsonData, err := json.MarshalIndent(elements, "", "    ")
	if err != nil {
		return err
	}

	_, err = file.Write(jsonData)
	if err != nil {
		return err
	}

	return nil
}

// Builds a CSV output for unprocessable files input
// I felt a CSV was easier to view for this purpose because the only
// header is the file name
func (c *client) CreateUnprocessableFilesCSV(files []string, fileName string) error {
	filePath := injectPath(fileName)
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	header := []string{"FileName"}
	if err := writer.Write(header); err != nil {
		return err
	}

	for _, file := range files {
		row := []string{
			file,
		}

		if err := writer.Write(row); err != nil {
			return err
		}
	}

	return nil
}

// Helper for placing output files into the results package
func injectPath(fileName string) string {
	filePath := filepath.Join("results", fileName)
	return filePath
}
