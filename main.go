package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/kikinteractive/go-geoip-service/service"
)

const IP_LOC = 7

func GetIps(csvPath string, ipChan chan<- string, endChan chan<- bool) {
	// Load a TXT file.
	f, _ := os.Open(csvPath)

	// Create a new reader.
	r := csv.NewReader(bufio.NewReader(f))
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		ip := record[IP_LOC]
		ipChan <- ip
	}
	endChan <- true
}

func ResolveIps(ipChan <-chan string, endChan <-chan bool) {

	for {
		select {
		case ip := <-ipChan:
			record, err := service.LookupIP(ip)
			if err == nil {
				regionCode := ""
				if record.RegionCode != nil {
					regionCode = *record.RegionCode
				}
				fmt.Printf("%v,%v,%v,%v,%v,%v,%v\n", ip, record.ContinentCode, record.CountryCode, regionCode, record.City, record.Location.Lat, record.Location.Lon)
			}
		case <-endChan:
			return
		}
	}
}

func main() {
	var dbPath string
	var csvPath string

	flag.StringVar(&dbPath, "db-path", "", "path to MaxMind GeoLite2 database")
	flag.StringVar(&csvPath, "csv-path", "", "path to csv file")
	flag.Parse()

	if 0 == len(dbPath) {
		log.Fatalln("you must specify a --db-path")
	}

	if 0 == len(csvPath) {
		log.Fatalln("you must specify a --csv-path")
	}

	service.LoadMaxmindDB(dbPath)

	ipChan := make(chan string)
	endChan := make(chan bool)
	go GetIps(csvPath, ipChan, endChan)
	ResolveIps(ipChan, endChan)
}
