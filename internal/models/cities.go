package models

// Defines a city object
type City struct {
	Latitude     string `json:"latitude"`
	Longitude    string `json:"longitude"`
	Geo          string `json:"geo"`
	City         string `json:"city"`
	ProvinceIcon string `json:"province_icon"`
	Province     string `json:"province"`
	CountryIcon  string `json:"country_icon"`
	Country      string `json:"country"`
}

// Defines a tmp package city object
type TmpCity struct {
	Latitude     string `json:"latitude"`
	Longitude    string `json:"longitude"`
	Geo          string `json:"geo"`
	City         string `json:"city"`
	ProvinceIcon string `json:"province_icon"`
	Province     string `json:"province"`
	CountryIcon  string `json:"country_icon"`
	Country      string `json:"country"`
	IsValid      bool   `json:"is_valid"`
}
