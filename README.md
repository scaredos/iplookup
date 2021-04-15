## IPLookup-UVPN
- Simple IPLookup API written for [UnknownVPN](https://unknownvpn.net/)'s use within the client.
- My public deployment of this API [ip.ddos.studio](https://ip.ddos.studio/v1/lookup/1.1.1.1)

- ( Requires GeoLite2 ASN and City database [.mmdb format] )
- ( Requires SSL Certificate ) 

## Usage
- `go run server.go`
- Request `http://0.0.0.0/v1/lookup/{IP}`
