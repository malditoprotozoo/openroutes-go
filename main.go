package main

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"strings"
)

var (
	Config Configuration
)

type OpenRoute struct {
	Key		string `yaml:"key"`
}

type Configuration struct {
	OpenRoute OpenRoute `yaml:"openroute"`
}

type OpenRoutesGeo struct {
	Geocoding struct {
		Version     string `json:"version"`
		Attribution string `json:"attribution"`
		Query       struct {
			Text            string `json:"text"`
			Size            int    `json:"size"`
			Private         bool   `json:"private"`
			BoundaryCountry string `json:"boundary.country"`
			Lang            struct {
				Name      string `json:"name"`
				Iso6391   string `json:"iso6391"`
				Iso6393   string `json:"iso6393"`
				Defaulted bool   `json:"defaulted"`
			} `json:"lang"`
			QuerySize int    `json:"querySize"`
			Parser    string `json:"parser"`
		} `json:"query"`
		Engine struct {
			Name    string `json:"name"`
			Author  string `json:"author"`
			Version string `json:"version"`
		} `json:"engine"`
		Timestamp int64 `json:"timestamp"`
	} `json:"geocoding"`
	Type     string `json:"type"`
	Features []struct {
		Type     string `json:"type"`
		Geometry struct {
			Type        string    `json:"type"`
			Coordinates []float64 `json:"coordinates"`
		} `json:"geometry"`
		Properties struct {
			ID           string  `json:"id"`
			Gid          string  `json:"gid"`
			Layer        string  `json:"layer"`
			Source       string  `json:"source"`
			SourceID     string  `json:"source_id"`
			Name         string  `json:"name"`
			Housenumber  string  `json:"housenumber"`
			Street       string  `json:"street"`
			Confidence   float64 `json:"confidence"`
			Accuracy     string  `json:"accuracy"`
			Country      string  `json:"country"`
			CountryGid   string  `json:"country_gid"`
			CountryA     string  `json:"country_a"`
			Region       string  `json:"region"`
			RegionGid    string  `json:"region_gid"`
			County       string  `json:"county"`
			CountyGid    string  `json:"county_gid"`
			Locality     string  `json:"locality"`
			LocalityGid  string  `json:"locality_gid"`
			LocalityA    string  `json:"locality_a"`
			Continent    string  `json:"continent"`
			ContinentGid string  `json:"continent_gid"`
			Label        string  `json:"label"`
		} `json:"properties"`
		Bbox []float64 `json:"bbox,omitempty"`
	} `json:"features"`
	Bbox []float64 `json:"bbox"`
}

func (c *Configuration) LoadConfig(filename string) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
		return
	} else {
		yaml.Unmarshal(bytes, &c)
	}
}


func (geo *OpenRoutesGeo) GeocodingPointOR(address, key string) (float64, float64, string, error) {
	var latitude, longitude float64
	var loc_type string
	countryCode := "CHL"
	address = strings.Replace(address, " ", "+", -1)
	url := "https://api.openrouteservice.org/geocode/search?text=" + address + "&boundary.country=" + countryCode
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", "Bearer " + key)
	req.Header.Add("Cache-Control", "no-cache")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	if err := json.Unmarshal(body, &geo); err != nil {
		return latitude, longitude, loc_type, err
	}
	if len(geo.Features) > 0 {
		latitude = geo.Features[0].Geometry.Coordinates[1]
		longitude = geo.Features[0].Geometry.Coordinates[0]
		loc_type = geo.Features[0].Properties.Accuracy

		return latitude, longitude, loc_type, nil
	} else {
		return latitude, longitude, loc_type, nil
	}
}

func main() {
	filename := "config.yaml"
	Config.LoadConfig(filename)
	query:= "santa blanca 1747"
	var geo OpenRoutesGeo
	fmt.Println(geo.GeocodingPointOR(query, Config.OpenRoute.Key))
}
