import logging
import requests 
from bs4 import BeautifulSoup

log_format = '%(asctime)s [%(name)s]  [%(levelname)s] : %(message)s' 
logging.basicConfig(level=logging.DEBUG, format=log_format)
logger = logging.getLogger('logger')

headers = { "User-agent": "i-do-not-care" }

req = requests.get('https://www.albumoftheyear.org', headers=headers) 
req.raise_for_status()

logger.info(f"{req.status_code = }")

soup = BeautifulSoup(req.content, "html.parser", from_encoding="utf_8")

content_container = soup.find('div', {"id": "centerContent"}) 

albums = content_container.find_all('div', {"class": "albumBlock"}) 
for album in albums[:5]: 
    title = album.find("div", {"class": "artistTitle"}).text 

    ratings_tags = album.find("div", {"class": "ratingRowContainer"}).find_all('div', {"class": "ratingRow"}) 
    ratings = {rating.find('div', {"class": "ratingText"}).text: rating.find('div', {"class": "ratingBlock"}).text for rating in ratings_tags} 

    logger.info(f"{title}, {ratings}")