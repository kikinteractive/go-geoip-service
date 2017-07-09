package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/kikinteractive/go-geoip-service/service"
)

func FloatToString(f float64) string {
	return strconv.FormatFloat(f, 'f', 4, 64)
}

func GetIps(input string, output string, ipLoc uint8) {
	fin, _ := os.Open(input)
	r := csv.NewReader(bufio.NewReader(fin))

	fout, _ := os.Create(output)
	w := csv.NewWriter(bufio.NewWriter(fout))
	defer w.Flush()

	for {
		row, err := r.Read()
		if err == io.EOF {
			break
		}
		ip := row[ipLoc]

		record, err := service.LookupIP(ip)
		if err == nil {
			regionCode := ""
			if record.RegionCode != nil {
				regionCode = *record.RegionCode
			}
			row = append(row, record.ContinentCode, record.CountryCode, regionCode, record.City, FloatToString(record.Location.Lat), FloatToString(record.Location.Lon))
		}
		w.Write(row)
	}
}

func main() {
	var dbPath string
	var csvPath string
	var outPath string

	flag.StringVar(&dbPath, "db-path", "", "path to MaxMind GeoLite2 database")
	flag.StringVar(&csvPath, "csv-path", "", "path to csv file")
	flag.StringVar(&outPath, "out-path", "", "path to out file")
	flag.Parse()

	if 0 == len(dbPath) {
		log.Fatalln("you must specify a --db-path")
	}

	if 0 == len(csvPath) {
		log.Fatalln("you must specify a --csv-path")
	}

	if 0 == len(outPath) {
		log.Fatalln("you must specify a --out-path")
	}

	service.LoadMaxmindDB(dbPath)

	GetIps(csvPath, outPath, 7)
}
