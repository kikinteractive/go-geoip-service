package service

import (
	"errors"
	"log"
	"net"
	"sync"

	geoip2 "github.com/oschwald/geoip2-golang"
)

var loaded bool
var mmdb *geoip2.Reader
var lock sync.Mutex

func init() {
	loaded = false
}

type Location struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

// Response is a struct that holds the data for the JSON HTTP response body.
type Response struct {
	CountryCode   string   `json:"country_code"`
	Country       string   `json:"country"`
	RegionCode    *string  `json:"region_code"`
	City          string   `json:"city"`
	ContinentCode string   `json:"continent_code"`
	Continent     string   `json:"continent"`
	Location      Location `json:"location"`
}

// LookupIP looks up the specified IP in the loaded Maxmind DB
func LookupIP(ip string) (*Response, error) {
	lock.Lock()
	defer lock.Unlock()

	if !loaded {
		return nil, errors.New("MaxMind DB not loaded")
	}

	parsedIP := net.ParseIP(ip) // nil result means error
	if nil == parsedIP {
		return nil, errors.New("failed to parse IP: " + ip)
	}

	record, err := mmdb.City(parsedIP)
	if err != nil {
		return nil, err
	}

	response := &Response{
		CountryCode:   record.Country.IsoCode,
		Country:       record.Country.Names["en"],
		City:          record.City.Names["en"],
		ContinentCode: record.Continent.Code,
		Continent:     record.Continent.Names["en"],
		Location: Location{
			Lat: record.Location.Latitude,
			Lon: record.Location.Longitude,
		},
	}

	if len(record.Subdivisions) > 0 {
		response.RegionCode = &record.Subdivisions[0].IsoCode
	}

	return response, nil
}

// LoadMaxmindDB loads a MaxMind DB into memory for use by the /lookup endpoint.
func LoadMaxmindDB(path string) {
	lock.Lock()
	defer lock.Unlock()

	if loaded {
		mmdb.Close()
		loaded = false
	}

	db, err := geoip2.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Loaded Maxmind DB from " + path)

	mmdb = db
	loaded = true
}

// UnloadMaxmindDB unloads the MaxMind DB from memory.  This is just for testing.
func UnloadMaxmindDB() {
	lock.Lock()
	defer lock.Unlock()

	if !loaded {
		return
	}

	log.Println("Unloaded MaxMind DB")

	mmdb.Close()
	loaded = false
}
