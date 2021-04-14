# copyright doink 2021
import sys
import maxminddb

if len(sys.argv) < 1:
    print(f"Usage: {sys.argv[0]} ip")
    exit(-1)

with maxminddb.open_database('GeoLite2-City.mmdb') as reader:
    r = reader.get(sys.argv[1])
    print(r['country']['iso_code'], "|", r['country']['names']['en'])
