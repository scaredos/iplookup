import ipsearch

# Initialize ipsearch client 
client = ipsearch.IPSearch()

# IP not specified (responds with the client's IP)
myip = client.lookup()
print(myip['ip'])

# Lookup IP
results = client.lookup('1.1.1.1')
print(results)
