package main

import (
	"fmt"
	"net"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"

	"github.com/caffix/cloudflare-roundtripper/cfrt"
	"github.com/gocolly/colly"
)

const (
	base_url    = "https://www.albumoftheyear.org"
	fetch_url   = "https://www.albumoftheyear.org/scripts/showMore.php"
	credits_url = "https://www.albumoftheyear.org/scripts/showAlbumCredits.php"
)

type Release struct {
	Artists    []string `json:"artists"`
	Featurings []string `json:"featurings"`
	Title      string   `json:"title"`
	Date       string   `json:"date"`
	Cover      string   `json:"cover"`
	Genres     []string `json:"genres"`
	Producers  []string `json:"producers"`
	Tracklist  []string `json:"tracklist"`
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
			getReleases("2023-04", 0, &allReleases[i], releaseType)
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
		link := e.DOM.Find("div.albumTitle").Parent().AttrOr("href", "")
		title := e.ChildText("div.albumTitle")

		release := Release{
			Artists:    []string{},
			Featurings: []string{},
			Title:      title,
			Date:       date,
			Cover:      cover,
			Genres:     []string{},
			Producers:  []string{},
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

	c1.OnHTML("div.artist", func(e *colly.HTMLElement) {
		e.ForEach("a", func(_ int, el *colly.HTMLElement) {
			details.Artists = append(details.Artists, el.Text)
		})
	})

	c1.OnHTML(`meta[itemprop="genre"]`, func(e *colly.HTMLElement) {
		genre := e.Attr("content")
		details.Genres = append(details.Genres, genre)
	})

	c1.OnHTML(".trackList ol li", func(e *colly.HTMLElement) {
		trackTitle := strings.TrimSpace(e.Text)
		details.Tracklist = append(details.Tracklist, trackTitle)

		titleContributors := parseTitle(trackTitle)

		if len(titleContributors[0]) > 0 {
			checkAddDuplicates(titleContributors[0], &details.Producers)
		}
		if len(titleContributors[1]) > 0 {
			checkAddDuplicates(titleContributors[1], &details.Featurings)
		}
	})

	c1.OnHTML("td.trackTitle", func(e *colly.HTMLElement) {
		if e.DOM.Find("span:first-child.featuredArtists").Length() > 0 {
			e.DOM.Find("span:first-child.featuredArtists").Remove()
		}

		var trackTitle strings.Builder
		var featuredArtists []string

		trackTitle.WriteString(e.DOM.Find("a:first-child").First().Text())

		e.ForEach("div.featuredArtists a", func(_ int, a *colly.HTMLElement) {
			featuredArtists = append(featuredArtists, a.Text)
		})

		if len(featuredArtists) > 0 {
			trackTitle.WriteString(" (feat. ")
			for i, artist := range featuredArtists {
				trackTitle.WriteString(artist)
				if i < len(featuredArtists)-1 {
					trackTitle.WriteString(", ")
				}
			}
			trackTitle.WriteString(")")
		}

		details.Tracklist = append(details.Tracklist, trackTitle.String())
		details.Featurings = append(details.Featurings, featuredArtists...)
	})

	c1.OnError(func(r *colly.Response, err error) {
		fmt.Println("ERROR:", r.StatusCode)
	})

	c1.Post(credits_url, data)
	c1.Visit(base_url + link)

	c1.Wait()
}

func checkAddDuplicates(candidates []string, contributors *[]string) {
	for _, name := range candidates {
		lowerName := strings.ToLower(name)
		lowerName = strings.Map(func(r rune) rune {
			if unicode.IsLetter(r) || unicode.IsDigit(r) {
				return r
			}
			return -1
		}, lowerName)

		found := false
		for _, contributor := range *contributors {
			lowerContributor := strings.ToLower(contributor)
			lowerContributor = strings.Map(func(r rune) rune {
				if unicode.IsLetter(r) || unicode.IsDigit(r) {
					return r
				}
				return -1
			}, lowerContributor)

			if strings.EqualFold(lowerName, lowerContributor) {
				found = true
				break
			}
		}
		if !found {
			*contributors = append(*contributors, name)
		}
	}
}

func parseTitle(songTitle string) [][]string {
	titleProducers := []string{}
	titleFeaturedArtists := []string{}

	songTitle = getTitleProducers(&songTitle, &titleProducers)

	re1 := regexp.MustCompile(`(?i)(?:(\[featuring|\(featuring| featuring)\.?|(\[feat|\(feat| feat)\.?|(\[ft|\(ft| ft)\.?)`)
	matchIndexes1 := re1.FindStringSubmatchIndex(songTitle)

	re2 := regexp.MustCompile(`(?i)(?:(\[with|\(with|\(w/|\[w/| w/)\.?)`)
	matchIndexes2 := re2.FindStringSubmatchIndex(songTitle)

	if len(matchIndexes1) != 0 && len(matchIndexes2) != 0 {
		if matchIndexes1[0] < matchIndexes2[0] {
			songTitle = getCoArtists(&songTitle, &titleFeaturedArtists, matchIndexes2)
			songTitle = getFeaturedArtists(&songTitle, &titleFeaturedArtists, matchIndexes1)
		} else {
			songTitle = getFeaturedArtists(&songTitle, &titleFeaturedArtists, matchIndexes1)
			songTitle = getCoArtists(&songTitle, &titleFeaturedArtists, matchIndexes2)
		}
	} else if len(matchIndexes1) != 0 {
		songTitle = getFeaturedArtists(&songTitle, &titleFeaturedArtists, matchIndexes1)
	} else if len(matchIndexes2) != 0 {
		songTitle = getCoArtists(&songTitle, &titleFeaturedArtists, matchIndexes2)
	}

	return [][]string{titleProducers, titleFeaturedArtists}
}

func getFeaturedArtists(songTitle *string, titleFeaturedArtists *[]string, matches []int) string {
	start := matches[0]
	end := matches[1]
	prodSubstring := ""
	if (*songTitle)[start] == '(' || (*songTitle)[start] == '[' {
		prodSubstring = cleanTitleSection((*songTitle)[start:])
		if (prodSubstring[0] == '(' && prodSubstring[len(prodSubstring)-1] == ')') || (prodSubstring[0] == '[' && prodSubstring[len(prodSubstring)-1] == ']') {
			prodSubstring = prodSubstring[(end - start) : len(prodSubstring)-1]
		}
	} else {
		prodSubstring = (*songTitle)[end:]
	}
	parseContributors(prodSubstring, titleFeaturedArtists)

	*songTitle = (*songTitle)[:start]

	return *songTitle
}

func getCoArtists(songTitle *string, titleFeaturedArtists *[]string, matches []int) string {
	start := matches[0]
	end := matches[1]
	prodSubstring := ""
	if (*songTitle)[start] == '(' || (*songTitle)[start] == '[' {
		prodSubstring = cleanTitleSection((*songTitle)[start:])
		if (prodSubstring[0] == '(' && prodSubstring[len(prodSubstring)-1] == ')') || (prodSubstring[0] == '[' && prodSubstring[len(prodSubstring)-1] == ']') {
			prodSubstring = prodSubstring[(end - start) : len(prodSubstring)-1]
		} else {
			prodSubstring = prodSubstring[(end - start):]
		}
	} else {
		prodSubstring = (*songTitle)[end:]
	}
	parseContributors(prodSubstring, titleFeaturedArtists)

	*songTitle = (*songTitle)[:start]

	return *songTitle
}

func getTitleProducers(songTitle *string, titleProducers *[]string) string {
	re := regexp.MustCompile(`(?i)(?:(\[prod. by|\(prod. by| prod. by|\[prod by|\(prod by| prod by|\[prod.|\(prod.| prod.|\[prod|\(prod| prod))`)
	matchesIndexes := re.FindAllStringSubmatchIndex(*songTitle, -1)

	if len(matchesIndexes) == 0 {
		return *songTitle
	}

	for _, match := range matchesIndexes {
		start := match[len(match)-2]
		end := match[len(match)-1]
		prodSubstring := ""
		if (*songTitle)[start] == '(' || (*songTitle)[start] == '[' {
			prodSubstring = cleanTitleSection((*songTitle)[start:])
			if (prodSubstring[0] == '(' && prodSubstring[len(prodSubstring)-1] == ')') || (prodSubstring[0] == '[' && prodSubstring[len(prodSubstring)-1] == ']') {
				prodSubstring = prodSubstring[(end - start) : len(prodSubstring)-1]
			}
		} else {
			prodSubstring = (*songTitle)[end:]
		}
		parseContributors(prodSubstring, titleProducers)
	}

	*songTitle = (*songTitle)[:matchesIndexes[0][0]]

	return *songTitle
}

func parseContributors(titleSection string, contributors *[]string) {
	separators := []string{"& ", ", ", "/ ", "\\ ", "feat. ", "ft. "}
	parts := splitString(titleSection, separators)
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" && part != " " && part != "." {
			part = cleanContributorName(part)
			*contributors = append(*contributors, part)
		}
	}
}

func splitString(str string, separators []string) []string {
	if len(separators) == 0 {
		return []string{str}
	}
	sep := separators[0]
	parts := strings.Split(str, sep)
	result := []string{}
	for _, part := range parts {
		result = append(result, splitString(part, separators[1:])...)
	}
	return result
}

func cleanTitleSection(str string) string {
	parenCount := 0
	bracketCount := 0
	m := len(str) - 1

	for i, c := range str {
		switch c {
		case '(':
			parenCount++
		case ')':
			parenCount--
		case '[':
			bracketCount++
		case ']':
			bracketCount--
		}

		if parenCount == 0 && bracketCount == 0 {
			m = i
			break
		}
	}

	prodSubstring := str[:m+1]

	return prodSubstring
}

func cleanContributorName(str string) string {

	lastChar := str[len(str)-1]

	if lastChar == ')' || lastChar == ']' {
		var charOpen rune
		var charClose rune

		if lastChar == ')' {
			charOpen = '('
			charClose = ')'
		} else {
			charOpen = '['
			charClose = ']'
		}

		count := 0

		for _, c := range str {
			switch c {
			case charOpen:
				count++
			case charClose:
				count--
			}
		}

		if count < 0 {
			str = str[:len(str)-1]
		}

	} else if lastChar == ',' {
		str = str[:len(str)-1]
	}

	return str
}
