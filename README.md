## IPLookup
- My public deployment of this API: [ip.ddos.studio](https://ip.ddos.studio/v1/lookup/1.1.1.1)

- ( Requires GeoLite2 ASN and City database [.mmdb format] )
- ( Requires SSL Certificate ) 

## Usage
- `go run httpserver.go`
- Request `http://0.0.0.0/v1/lookup/{IP}`

## Documentation
- Sample Request: `GET /v1/lookup/1.1.1.1` OR `GET /v1/ip`
- Sample Response:
```
{
    "ip": "1.1.1.1", 
    "hostname": "one.one.one.one.", 
    "latitude": -33.494, 
    "longitude": 143.2104,
    "country": "Australia",
    "org": "CLOUDFLARENET",
    "asn": 13335,
    "timezone": "Australia/Sydney",
    "country_code": "AU",
    "registered": true
}
```

- The `registered` flag displays `true` when the IP has a corresponding organization and hostname
