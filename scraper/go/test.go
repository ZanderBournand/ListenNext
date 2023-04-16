package main

import (
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/caffix/cloudflare-roundtripper/cfrt"
	"github.com/gocolly/colly"
)

const (
	base_url    = "https://www.albumoftheyear.org"
	fetch_url   = "https://www.albumoftheyear.org/scripts/showMore.php"
	credits_url = "https://www.albumoftheyear.org/scripts/showAlbumCredits.php"
)

type Release struct {
	Artist    string
	Title     string
	Date      string
	Cover     string
	Genres    []string
	Producers []string
}

func main() {
	releaseTypes := []string{"lp", "ep", "single", "mixtape", "reissue"}
	allReleases := make([][]Release, len(releaseTypes))

	var wg sync.WaitGroup
	wg.Add(len(releaseTypes))

	for i, releaseType := range releaseTypes {
		fmt.Println("Fetching " + releaseType + "s...")
		go func(i int, releaseType string) {
			defer wg.Done()
			getReleases("2023-03", 0, &allReleases[i], releaseType)
		}(i, releaseType)
	}

	wg.Wait()

	for i, releaseType := range releaseTypes {
		fmt.Println(releaseType, "count:", len(allReleases[i]))
	}

}

func getReleases(date string, start int, allReleases *[]Release, releaseType string) {
	data := map[string]string{
		"type":      "albumMonth",
		"sort":      "release",
		"albumType": releaseType,
		"start":     fmt.Sprintf("%v", start),
		"date":      date,
		"genre":     "",
		"reviews":   "",
	}

	client := &http.Client{
		Timeout: 15 * time.Second,
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   15 * time.Second,
				KeepAlive: 15 * time.Second,
				DualStack: true,
			}).DialContext,
		},
	}
	client.Transport, _ = cfrt.New(client.Transport)

	c := colly.NewCollector()
	c.WithTransport(client.Transport)

	count := 0
	var wg sync.WaitGroup

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("ERROR:", r.StatusCode)
	})

	c.OnHTML("div.albumBlock", func(e *colly.HTMLElement) {
		date := e.ChildText("div.date")
		cover := e.ChildAttr("img.lazyload", "data-src")
		artist := e.ChildText("div.artistTitle")
		link := e.DOM.Find("div.albumTitle").Parent().AttrOr("href", "")
		title := e.ChildText("div.albumTitle")

		release := Release{
			Artist:    artist,
			Title:     title,
			Date:      date,
			Cover:     cover,
			Genres:    []string{},
			Producers: []string{},
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			getDetails(link, &release)
			*allReleases = append(*allReleases, release)
		}()

		count++

		if count == 60 {
			count = 0
			newStart, _ := strconv.Atoi(data["start"])
			newStart += 60
			data["start"] = strconv.Itoa(newStart)
			c.Post(fetch_url, data)
		}
	})

	c.Post(fetch_url, data)
	c.Wait()
	wg.Wait()
}

func getDetails(link string, details *Release) {
	id := strings.Split(strings.Split(link, "/")[2], "-")[0]

	data := map[string]string{
		"albumID": id,
	}

	client := &http.Client{
		Timeout: 15 * time.Second,
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   15 * time.Second,
				KeepAlive: 15 * time.Second,
				DualStack: true,
			}).DialContext,
		},
	}
	client.Transport, _ = cfrt.New(client.Transport)

	c1 := colly.NewCollector()
	c1.WithTransport(client.Transport)

	c1.OnHTML("td.name a", func(e *colly.HTMLElement) {
		details.Producers = append(details.Producers, e.Text)
	})
	c1.OnHTML("a[itemprop=genre]", func(e *colly.HTMLElement) {
		details.Genres = append(details.Genres, e.Text)
	})

	c1.OnError(func(r *colly.Response, err error) {
		fmt.Println("ERROR:", r.StatusCode)
	})

	c1.Post(credits_url, data)
	c1.Visit(base_url + link)

	c1.Wait()
}
