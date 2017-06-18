package main

import (
	"encoding/json"
	"flag"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/kikinteractive/go-geoip-service/service"
)

type errorResponse struct {
	Error string `json:"error"`
}

type IpList struct {
	Ips []string `json:"ips"`
}

func writeErrorResponse(err error, w http.ResponseWriter) {
	resp := errorResponse{Error: err.Error()}

	bytes, jsonErr := json.Marshal(resp)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	w.WriteHeader(http.StatusInternalServerError)
	io.WriteString(w, string(bytes))
}

func lookupHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ip := req.URL.Query().Get("ip")

	record, err := service.LookupIP(ip)
	if err != nil {
		writeErrorResponse(err, w)
		return
	}

	bytes, err := json.Marshal(record)
	if err != nil {
		log.Fatal(err)
	}
	io.WriteString(w, string(bytes))
}

func multiLookupHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	decoder := json.NewDecoder(req.Body)
	defer req.Body.Close()

	var ips IpList
	err := decoder.Decode(&ips)
	if err != nil {
		writeErrorResponse(err, w)
		return
	}

	record, err := service.MultiLookupIP(ips.Ips)
	if err != nil {
		writeErrorResponse(err, w)
		return
	}

	bytes, err := json.Marshal(record)
	if err != nil {
		log.Fatal(err)
	}
	io.WriteString(w, string(bytes))
}

func main() {
	var dbPath string
	var port int

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
	http.HandleFunc("/multi-lookup", multiLookupHandler)
	log.Fatal(http.ListenAndServe(":"+stringPort, nil))
}
