// ipsearch-api | version 0.1 (dev)
// dev
// https://download.maxmind.com/app/geoip_download?edition_id=GeoLite2-City&license_key=eOFeNYpRl2MimT1i&suffix=tar.gz
package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os/exec"
	"regexp"
	"strings"
)

func ipClean(ip string) string {
	r, _ := regexp.Compile(`\b(([0-9a-fA-F]{1,4}:){7,7}[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,7}:|([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})|:((:[0-9a-fA-F]{1,4}){1,7}|:)|fe80:(:[0-9a-fA-F]{0,4}){0,4}%[0-9a-zA-Z]{1,}|::(ffff(:0{1,4}){0,1}:){0,1}((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])|([0-9a-fA-F]{1,4}:){1,4}:((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9]))\b|\b(?:[0-9]{1,3}\.){3}[0-9]{1,3}\b`)
	return r.FindString(ip)
}

func ipLookup(w http.ResponseWriter, r *http.Request) {
	// Response headers
	w.Header().Set("Cache-Control", "public, max-age=2678400, s-max-age=2678400")
	w.Header().Set("Content-Type", "application/json")
	ip := ipClean(r.URL.Path)
	if ip == "" {
		fmt.Fprintf(w, "{\"error\": \"no ip provided (/v1/lookup/{ip})\"}\n")
		return
	}
	//cmd, err := exec.Command("python3", "max.py", ip).Output()
	cmd, err := exec.Command("python", "max.py", ip).Output()
	if err != nil {
		fmt.Fprintf(w, "{\"error\": \"invalid ip provided\"}\n")
		return
	}
	out := strings.Split(string(cmd), "|")        // split at |
	out[0] = strings.Replace(out[0], " ", "", 1)  // remove spaces
	out[1] = strings.TrimSuffix(out[1][1:], " ")  // remove spaces
	out[2] = strings.Replace(out[2], " ", "", 1)  // remove spaces
	out[3] = strings.Replace(out[3], " ", "", 1)  // remove spaces
	out[9] = strings.Replace(out[9], "\n", "", 1) // Remove newline char
	ips := make(map[string]string)
	ips["ip"] = ip
	hosts, errz := net.LookupAddr(ip)
	if errz != nil {
		ips["hostname"] = ""
	} else {
		ips["hostname"] = string(hosts[0])
	}
	ips["country"] = strings.Replace(out[1], " ", "_", -1)
	ips["country_spaces"] = out[1]
	ips["continent"] = strings.TrimSuffix(out[5][1:], " ")
	ips["continent_code"] = out[9][1:]
	ips["latitude"] = strings.TrimSuffix(out[6][1:], " ")
	ips["longitude"] = strings.TrimSuffix(out[7][1:], " ")
	ips["registered"] = strings.TrimSuffix(out[8][1:], " ")
	ips["timezone"] = strings.TrimSuffix(out[4][1:], " ")
	ips["country_code"] = out[0]
	ips["asn"] = strings.TrimSuffix(out[2], " ")
	ips["org"] = strings.TrimSuffix(out[3], " ")
	data, _ := json.MarshalIndent(ips, "", "    ")
	fmt.Fprintf(w, string(data))
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
	w.Header().Set("cache-control", "public, max-age=2678400, s-max-age=2678400")
	w.Header().Set("Content-Type", "application/json")
	ip := ipClean(string(r.Header.Get("Cf-Connecting-Ip")))
	//cmd, err := exec.Command("python3", "max.py", ip).Output()
	cmd, err := exec.Command("python", "max.py", ip).Output()
	if err != nil {
		fmt.Fprintf(w, "{\"error\": \"api does not support ipv6\"}\n")
		return
	}

	out := strings.Split(string(cmd), "|")        // split at |
	out[0] = strings.Replace(out[0], " ", "", 1)  // remove spaces
	out[1] = strings.TrimSuffix(out[1][1:], " ")  // remove spaces
	out[2] = strings.Replace(out[2], " ", "", 1)  // remove spaces
	out[3] = strings.Replace(out[3], " ", "", 1)  // remove spaces
	out[9] = strings.Replace(out[9], "\n", "", 1) // Remove newline char
	ips := make(map[string]string)
	ips["ip"] = ip
	hosts, errz := net.LookupAddr(ip)
	if errz != nil {
		ips["hostname"] = ""
	} else {
		ips["hostname"] = string(hosts[0])
	}
	ips["country"] = strings.Replace(out[1], " ", "_", -1)
	ips["country_spaces"] = out[1]
	ips["continent"] = strings.TrimSuffix(out[5][1:], " ")
	ips["continent_code"] = out[9][1:]
	ips["latitude"] = strings.TrimSuffix(out[6][1:], " ")
	ips["longitude"] = strings.TrimSuffix(out[7][1:], " ")
	ips["registered"] = strings.TrimSuffix(out[8][1:], " ")
	ips["timezone"] = strings.TrimSuffix(out[4][1:], " ")
	ips["country_code"] = out[0]
	ips["asn"] = strings.TrimSuffix(out[2], " ")
	ips["org"] = strings.TrimSuffix(out[3], " ")
	data, _ := json.MarshalIndent(ips, "", "    ")
	fmt.Fprintf(w, string(data))
}
func main() {
	fmt.Println("webserver started")
	http.HandleFunc("/v1/ip", ipRes)
	http.HandleFunc("/v1/lookup/", ipLookup)
	http.HandleFunc("/", getStarted)
	http.ListenAndServe("0.0.0.0:80", nil)
	//_ = http.ListenAndServeTLS("0.0.0.0:443", "origin.pem", "private.pem", nil)
}
