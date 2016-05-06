package handlers

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"

	"github.com/KSubedi/ip-geolocation-api-server/utils"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/oschwald/geoip2-golang"
)

// ReturnData struct will hold information that will later be returned back to
type ReturnDataV1 struct {
	IP         string  `json:"ip"`
	City       string  `json:"city"`
	Region     string  `json:"region"`
	Country    string  `json:"country"`
	Postal     string  `json:"postal"`
	Continent  string  `json:"continent"`
	HostName   string  `json:"hostname"`
	SearchName string  `json:"searchname"`
	Latitude   float64 `json:"latitude"`
	Longitude  float64 `json:"longitude"`
}

// If there is an error message, return it here
type MessageDataV1 struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
}

func V1(w http.ResponseWriter, r *http.Request) {
	db := utils.DB
	var ipObj net.IP
	searchname := ""
	queryIP := mux.Vars(r)["ip"]
	if queryIP == "" {
		// Grab remote IP Address & Port
		ipPort := r.RemoteAddr

		// Extract IP from the IP & Port
		ipString, _, err := net.SplitHostPort(ipPort)
		if err != nil {
			fmt.Println(err)
		}

		// Try to parse IP and make sure it is valid
		ipObj = net.ParseIP(ipString)
		if ipObj == nil {
			eMessage := MessageDataV1{true, "Could Not Parse IP!"}
			encoder := json.NewEncoder(w)
			encoder.Encode(eMessage)

			logs := logrus.Fields{
				"error": "true",
				"type":  "parse_ip_error",
				"ip":    "ipstring",
			}
			utils.LogError(logs, "Could Not Parse IP")

			return
		}
	} else {
		// ipString := strings.TrimPrefix(r.URL.Path, "/v1/")
		ipString := queryIP

		// Try to parse IP and make sure it is valid
		ipObj = net.ParseIP(ipString)
		if ipObj == nil { // If the IP is not valid, try to parse as hostname
			ips, err := net.LookupIP(ipString)

			if err != nil || len(ips) == 0 {
				eMessage := MessageDataV1{true, "IP or Hostname Invalid!"}
				encoder := json.NewEncoder(w)
				encoder.Encode(eMessage)

				logs := logrus.Fields{
					"error": "true",
					"type":  "invalid_hostname",
					"ip":    "ipstring",
				}
				utils.LogError(logs, "IP or Hostname Invalid")

				return
			} else {
				searchname = ipString
				ipObj = ips[0]
			}

		}
	}

	returnIP(ipObj, searchname, r, w, db)

}

// Return the IP after processing
func returnIP(ip net.IP, searchname string, r *http.Request, w http.ResponseWriter, db *geoip2.Reader) {
	// Set content type to JSON
	w.Header().Set("Content-Type", "application/json")

	// Find the record
	record, err := db.City(ip)
	if err != nil {
		fmt.Println(err)
	}

	// Grab data out of the record
	city := record.City.Names["en"]
	country := record.Country.Names["en"]
	lat := record.Location.Latitude
	long := record.Location.Longitude
	post := record.Postal.Code
	continent := record.Continent.Names["en"]
	hostname := ""

	// Lookup hostname if not provided
	names, err := net.LookupAddr(ip.String())
	if err != nil {
		fmt.Println(err)
	}

	if len(names) > 0 {
		hostname = names[0]
	}

	region := ""
	// Loop through region array
	for index, division := range record.Subdivisions {
		if index == 0 {
			region += division.Names["en"]
		} else {
			region += ", " + division.Names["en"]
		}
	}

	requestIP := r.RemoteAddr
	// Extract IP from the IP & Port
	ipString, _, err := net.SplitHostPort(requestIP)
	if err == nil {
		requestIP = ipString
	}

	// Log things properly
	logs := logrus.Fields{
		"ip":          ip.String(),
		"city":        city,
		"region":      region,
		"country":     country,
		"postal":      post,
		"continent":   continent,
		"hostname":    hostname,
		"searchname":  searchname,
		"lat":         lat,
		"long":        long,
		"remote-host": requestIP,
		"path":        r.RequestURI,
	}
	msg := "Manual Search"
	if r.RequestURI == "/v1/" {
		msg = "Self Search"
	}

	go utils.LogInfo(logs, msg)

	// Create a struct to return & return after encoding it to JSon
	returnData := ReturnDataV1{ip.String(), city, region, country, post, continent, hostname, searchname, lat, long}
	jsonEncoder := json.NewEncoder(w)
	jsonEncoder.Encode(returnData)

}
