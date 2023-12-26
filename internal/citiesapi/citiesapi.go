package citiesapi

import (
	"github.com/efuchsman/city_validation_s3/internal/models"
	log "github.com/sirupsen/logrus"
)

type Client interface {
	ReturnValidTmpElements() ([]*models.TmpCity, error)
	ReturnInvalidTmpElements() ([]*models.TmpCity, error)
	ReturnUnprocessableFiles() *UnprocessableFiles
	CreateValidElementsJSON(elements []*models.TmpCity, fileName string) error
	CreateInvalidElementsJSON(elements []*models.TmpCity, fileName string) error
	CreateUnprocessableFilesCSV(files []string, fileName string) error
	ReturnCitiesMap() *CitiesMap
}

type client struct {
	tmpCities []*models.TmpCity
	badFiles  *UnprocessableFiles
	cities    *CitiesMap
}

func NewClient() Client {
	tmp, bad, err := BuildTMPCities()
	if err != nil {
		log.Fatalf("Error building cities from tmp files: %v", err)
	}

	citiesMap, err := BuildCitiesMapFromCities()
	if err != nil {
		log.Fatalf("Error building cities from cities map: %v", err)
	}

	return &client{
		tmpCities: tmp,
		badFiles:  bad,
		cities:    citiesMap,
	}
}
