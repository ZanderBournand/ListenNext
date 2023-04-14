
import json
from urllib import parse
from urllib.request import Request, urlopen
from bs4 import BeautifulSoup
import requests

url = 'https://www.albumoftheyear.org/scripts/showMore.php'

data = {
    'type': 'albumMonth',
    'sort': '',
    'albumType': 'single',
    'start': '0',
    'date': '2023-04',
    'genre': '',
    'reviews': '',
}

data = parse.urlencode(data).encode('utf-8')

headers = {
    'User-Agent': 'Mozilla/6.0',
}

req = Request(url, data=data, headers=headers)
response = urlopen(req)

print(response.status)