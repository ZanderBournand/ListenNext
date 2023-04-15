import requests
from bs4 import BeautifulSoup

base_url = 'https://www.albumoftheyear.org'
fetch_url = 'https://www.albumoftheyear.org/scripts/showMore.php'
credits_url = 'https://www.albumoftheyear.org/scripts/showAlbumCredits.php'

headers = {
    "User-Agent": "Mozilla/6.0",
}


def getDetails(link, details):
    id = link.split('/')[2].split('-')[0]

    data = {
        'albumID': id,
    }

    # Getting producer credits
    res = requests.post(credits_url, headers=headers, data=data)
    soup = BeautifulSoup(res.content, "html.parser", from_encoding="utf_8")

    producers = soup.find_all('td', {"class": "name"})
    for producer in producers:
        details["producers"].append(producer.a.text)

    # Getting album/song general page
    res = requests.get(base_url + link, headers=headers)
    soup = BeautifulSoup(res.content, "html.parser", from_encoding="utf_8")

    # Getting genres of music
    genres = soup.find_all('a', {"itemprop": "genre"})
    for genre in genres:
        details["genres"].append(genre.text)


def getReleases(date, start, all_releases, type):
    data = {
        'type': 'albumMonth',
        'sort': 'release',
        'albumType': type,
        'start': start,
        'date': date,
        'genre': '',
        'reviews': '',
    }

    res = requests.post(fetch_url, headers=headers, data=data)

    soup = BeautifulSoup(res.content, "html.parser", from_encoding="utf_8")
    releases = soup.find_all('div', {"class": "albumBlock"})

    count = 0
    for release in releases:
        date = release.find("div", {"class": "date"}).text
        img_tag = release.find("img", {"class": "lazyload"})
        if img_tag:
            cover = img_tag.get("data-src")
        else:
            cover = None
        artist = release.find("div", {"class": "artistTitle"}).text
        link = release.find("div", {"class": "albumTitle"}).parent.get("href")
        title = release.find("div", {"class": "albumTitle"}).text

        count += 1
        release_dict = {
            "artist": artist,
            "title": title,
            "date": date,
            "cover": cover,
            "genres": [],
            "producers": [],
        }

        getDetails(link, release_dict)
        all_releases.append(release_dict)

    if (count == 60):
        getReleases('2023-04', start+60, all_releases, type)


albums = []
eps = []
songs = []

print("Fetching April Albums...")
getReleases('2023-04', 0, albums, 'lp')

print("Fetching April EPs...")
getReleases('2023-04', 0, eps, 'ep')

print("Fetching April Songs...")
getReleases('2023-04', 0, songs, 'single')

print("Albums:", len(albums))
print("EPs:", len(eps))
print("Songs:", len(songs))
