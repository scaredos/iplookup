# Simple API wrapper for this iplookup
# Copyright (c) 2021 scaredos
import requests


class IPLookup:
    def __init__(self):
        self.base_url = 'https://ip.ddos.studio/v1'

    def lookup(self, ip: str = None):
        """
        Lookup an IP address with the publically available API (ip.ddos.studio)

        :param ip: The IP address to search/lookup
        :return: A JSON object with the results of the search
        """

        if ip == None:
            response = requests.get(self.base_url + '/ip')
            return response.json()

        response = requests.get(self.base_url + '/lookup/' + ip)

        if response.status_code != 200:
            raise LookupError(
                "Unable to lookup IP: Server down or client IP address blocked")
        elif "error" in response.json():
            raise LookupError(
                "Unable to lookup IP: IP address provided is invalid")
        else:
            return response.json()
