package main

import (
	"encoding/json"
	"fmt"
	"github.com/oschwald/maxminddb-golang"
	"net"
	"net/http"
	"regexp"
)

type bogon struct {
	Ip         string `json:"ip"`
	Bogon      bool   `json:"bogon"`
	Registered bool   `json:"registered"`
}

type ipResponse struct {
	Ip          string  `json:"ip"`
	Hostname    string  `json:"hostname"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Country     string  `json:"country"`
	Org         string  `json:"org"`
	Asn         uint64  `json:"asn"`
	Timezone    string  `json:"timezone"`
	CountryCode string  `json:"country_code"`
	Registered  bool    `json:"registered"`
}

type asRecord struct {
	autonomous_system_orginization string `maxminddb:"autonomous_system_orginization"`
	autonomous_system_number       int    `maxminddb:"autonomous_system_number"`
}

type Record struct {
	Country struct {
		IsoCode string `maxminddb:"iso_code"`
	}
	Location struct {
		Latitude  float64 `maxminddb:"latitude"`
		Longitude float64 `maxminddb:"longitude"`
		TimeZone  string  `maxminddb:"time_zone"`
	}
}

func ipClean(ip string) string {
	// Regex for an IPv4 or IPv6 address
	r, _ := regexp.Compile(`\b(([0-9a-fA-F]{1,4}:){7,7}[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,7}:|([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})|:((:[0-9a-fA-F]{1,4}){1,7}|:)|fe80:(:[0-9a-fA-F]{0,4}){0,4}%[0-9a-zA-Z]{1,}|::(ffff(:0{1,4}){0,1}:){0,1}((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])|([0-9a-fA-F]{1,4}:){1,4}:((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9]))\b|\b(?:[0-9]{1,3}\.){3}[0-9]{1,3}\b`)
	return r.FindString(ip)
}

func ipLookup(w http.ResponseWriter, r *http.Request) {
	// Response headers
	w.Header().Set("Cache-Control", "public, max-age=2678400, s-max-age=2678400")
	w.Header().Set("Content-Type", "application/json")

	// Clean the user input to prevent any exploitation
	ip := ipClean(r.URL.Path)

	// If there is no IP in the URL Path, return error
	if ip == "" {
		fmt.Fprintf(w, "{\"error\": \"no ip provided (/v1/lookup/{ip})\"}\n")
		return
	}

	// Open database and search for IP address
	db, _ := maxminddb.Open("GeoLite2-City.mmdb")
	var record map[string]map[string]interface{}
	err := db.Lookup(net.ParseIP(ip), &record)
	// Check if IP address in database
	if len(record) == 0 {
		var bogonRes bogon
		bogonRes.Ip = ip
		bogonRes.Bogon = true
		bogonRes.Registered = false
		re, _ := json.MarshalIndent(bogonRes, "", "    ")
		fmt.Fprintf(w, string(re))
		return
	}
	var resp ipResponse
	resp.Ip = ip
	hosts, err := net.LookupAddr(ip)
	if err != nil {
		resp.Hostname = ""
		fmt.Println(err)
	} else {
		resp.Hostname = hosts[0]
		resp.Registered = true
	}
	resp.Country = record["country"]["names"].(map[string]interface{})["en"].(string)
	resp.Latitude = record["location"]["latitude"].(float64)
	resp.Longitude = record["location"]["longitude"].(float64)
	resp.Timezone = record["location"]["time_zone"].(string)
	resp.CountryCode = record["country"]["iso_code"].(string)
	defer db.Close()
	db, _ = maxminddb.Open("GeoLite2-ASN.mmdb")
	var asnrecord map[string]interface{}
	err = db.Lookup(net.ParseIP(ip), &asnrecord)
	// Check if IP address in database
	if len(asnrecord) == 0 {
		var bogonRes bogon
		bogonRes.Ip = ip
		bogonRes.Bogon = true
		bogonRes.Registered = false
		re, _ := json.MarshalIndent(bogonRes, "", "    ")
		fmt.Fprintf(w, string(re))
		return
	}
	if asnrecord["autonomous_system_number"] == nil {
		resp.Asn = 0
		resp.Org = ""
		resp.Registered = false
	} else {
		resp.Asn = asnrecord["autonomous_system_number"].(uint64)
		resp.Org = asnrecord["autonomous_system_organization"].(string)
		resp.Registered = true
	}
	re, _ := json.MarshalIndent(resp, "", "    ")
	fmt.Fprintf(w, string(re))
}

func getStarted(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("cache-control", "public, max-age=2678400, s-max-age=2678400")
	w.Header().Set("Content-Type", "application/json")
	res := make(map[string]map[string]string)
	res["1.info"] = make(map[string]string)
	res["1.info"]["1.url"] = "/v1/lookup/:ip"
	res["1.info"]["2.info"] = "respond with info on ip address"
	res["1.info"]["3.headers"] = "no headers required"
	res["2.info"] = make(map[string]string)
	res["2.info"]["1.uri"] = "/v1/ip"
	res["2.info"]["2.info"] = "respond with client ip address and info"
	res["2.info"]["3.headers"] = "no headers required"
	data, _ := json.MarshalIndent(res, "", "    ")
	fmt.Fprintf(w, string(data))
}

func ipRes(w http.ResponseWriter, r *http.Request) {
	ip := ipClean(string(r.Header.Get("Cf-Connecting-Ip")))
	w.Header().Set("Cache-Control", "public, max-age=2678400, s-max-age=2678400")
	w.Header().Set("Content-Type", "application/json")
	db, _ := maxminddb.Open("GeoLite2-City.mmdb")
	var record map[string]map[string]interface{}
	err := db.Lookup(net.ParseIP(ip), &record)
	if err != nil {
		fmt.Println(err)
	}
	var resp ipResponse
	resp.Ip = ip
	hosts, err := net.LookupAddr(ip)
	if err != nil {
		resp.Hostname = ""
		resp.Registered = false
		fmt.Println(err)
	} else {
		resp.Hostname = hosts[0]
		resp.Registered = true
	}
	resp.Country = record["country"]["names"].(map[string]interface{})["en"].(string)
	resp.Latitude = record["location"]["latitude"].(float64)
	resp.Longitude = record["location"]["longitude"].(float64)
	resp.Timezone = record["location"]["time_zone"].(string)
	resp.CountryCode = record["country"]["iso_code"].(string)
	defer db.Close()
	db, _ = maxminddb.Open("GeoLite2-ASN.mmdb")
	var asnrecord map[string]interface{}
	err = db.Lookup(net.ParseIP(ip), &asnrecord)
	if err != nil {
		fmt.Println(err)
	} else {
		resp.Asn = asnrecord["autonomous_system_number"].(uint64)
		resp.Org = asnrecord["autonomous_system_organization"].(string)
	}
	re, _ := json.MarshalIndent(resp, "", "    ")
	fmt.Fprintf(w, string(re))
}

func main() {
	fmt.Println("webserver started")
	http.HandleFunc("/v1/ip", ipRes)
	http.HandleFunc("/v1/lookup/", ipLookup)
	http.HandleFunc("/", getStarted)
	http.ListenAndServe("0.0.0.0:3600", nil)
}
