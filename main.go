package main

import (
	"encoding/json"
	"flag"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/dseevr/go-geoip-service/service"
)

func lookupHandler(w http.ResponseWriter, req *http.Request) {
	ip := req.URL.Query().Get("ip")

	record, _ := service.LookupIP(ip)

	bytes, err := json.Marshal(record)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, string(bytes))
}

// -------------------------------------------------------------------------------------------------

var (
	dbPath string
	port   int
)

func main() {
	flag.StringVar(&dbPath, "db-path", "", "path to MaxMind GeoLite2 database")
	flag.IntVar(&port, "port", 12345, "http port to listen on")
	flag.Parse()

	if 0 == len(dbPath) {
		log.Fatalln("you must specify a --db-path")
	}

	// TODO: allow port 0? not sure if it's worth it
	if port < 1 || port > 65535 {
		log.Fatalln("--port must be >= 1 and <= 65535")
	}

	service.LoadMaxmindDB(dbPath)

	stringPort := strconv.Itoa(port)

	log.Println("Listening on 0.0.0.0:" + stringPort)

	http.HandleFunc("/lookup", lookupHandler)
	log.Fatal(http.ListenAndServe(":"+stringPort, nil))
}
