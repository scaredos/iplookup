## IPLookup
- A simple and fast IP Information API written in Go
- There is a public deployment of this API available at [ip.ddos.studio](https://ip.ddos.studio/v1/lookup/1.1.1.1)

## Prerequisites
- **You must have GeoLite2 databases**
- Install all dependencies
    - Install Go
    - `go get github.com/oschwald/maxminddb-golang`

## How to run
- You can run it via the following commands:
    - `go build httpserver.go && ./httpserver` OR
    - `go run httpserver.go`


## Features
- lookup
   - Returns information about the provided IP address
   - Example response available [here](https://ip.ddos.studio/v1/lookup/1.1.1.1)
    
-  ip
    - Returns information about the client's IP address
    - Same response as `lookup`


## Dependencies
- Go
- GeoLite2 Databases
- [github.com/oschwald/maxminddb-golang](https://github.com/oschwald/maxminddb-golang)
