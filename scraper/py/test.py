import requests

url = 'https://www.albumoftheyear.org/scripts/showMore.php'

headers = {
 "User-Agent": "Mozilla/6.0",
}

data = {
    'type': 'albumMonth',
    'sort': '',
    'albumType': 'single',
    'start': '0',
    'date': '2023-04',
    'genre': '',
    'reviews': '',
}

# Second request using requests -> Forbidden 403
res = requests.post(url, headers=headers, data=data)
print("STATUS CODE:", res.status_code)
print(res.text)
