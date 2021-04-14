// ipsearch-api | version 0.1 (dev)
// dev
package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os/exec"
	"strings"
)

func ipLookup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("cache-control", "public")
	w.Header().Set("Content-Type", "application/json")
	uri := strings.Split(r.URL.Path, "/")
	if len(uri) < 2 {
		fmt.Fprintf(w, "{\"Error\": \"invalid ip provided\"}\n")
		return
	}
	ip := string(uri[len(uri)-1])
	// lets clean up this ip for rce vulns
	ip = strings.Replace(ip, "&", "", -1)
	ip = strings.Replace(ip, "|", "", -1)
	ip = strings.Replace(ip, ";", "", -1)
	ip = strings.Replace(ip, " ", "", -1)
	ip = strings.Replace(ip, "/", "", -1)
	ip = strings.Replace(ip, "\\", "", -1)
	ip = strings.Replace(ip, "-", "", -1)
	// we have to read from a python since maxminddb golang api not official

	// dependant upon your OS, switch the following statements
	cmd, err := exec.Command("python3", "max.py", ip).Output()
	//cmd, err := exec.Command("python", "max.py", ip).Output()
	if err != nil {
		fmt.Fprintf(w, "{\"Error\": \"invalid ip provided\"}\n")
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
	http.ListenAndServe("0.0.0.0:2096", nil)
}
