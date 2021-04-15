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
	r, _ := regexp.Compile(`\b(?:[0-9]{1,3}\.){3}[0-9]{1,3}\b`)
	return r.FindString(ip)
}

func ipLookup(w http.ResponseWriter, r *http.Request) {
	// response headers for browsers
	w.Header().Set("cache-control", "public, max-age=2678400, s-max-age=2678400")
	w.Header().Set("Content-Type", "application/json")
	uri := strings.Split(r.URL.Path, "/")
	if len(uri) < 2 {
		fmt.Fprintf(w, "{\"error\": \"no ip provided (/v1/lookup/{ip})\"}\n")
		return
	}
	ip := string(uri[len(uri)-1])
	// clean user input and return ip addresses
	ip = ipClean(ip)
	cmd, err := exec.Command("python3", "max.py", ip).Output()
	//cmd, err := exec.Command("python", "max.py", ip).Output()
	if err != nil {
		fmt.Fprintf(w, "{\"error\": \"invalid ip provided\"}\n")
		fmt.Println(string(cmd))
		fmt.Println(err)
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
	ips["asn"] = strings.TrimSuffix(out[2][1:], " ")
	ips["org"] = strings.TrimSuffix(out[3], " ")
	data, _ := json.MarshalIndent(ips, "", "    ")
	fmt.Fprintf(w, string(data))
}

func main() {
	fmt.Println("webserver started")
	http.HandleFunc("/v1/lookup/", ipLookup)
	http.ListenAndServe("0.0.0.0:80", nil)
}
