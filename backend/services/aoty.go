package services

import (
	"errors"
	"fmt"
	"main/db"
	"main/types"
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
	base_url          = "https://www.albumoftheyear.org"
	fetch_url         = "https://www.albumoftheyear.org/scripts/showMore.php"
	credits_url       = "https://www.albumoftheyear.org/scripts/showAlbumCredits.php"
	parseDateFormat   = "Jan 2 2006"
	requestDateFormat = "2006-01"
)

func Upload(releases map[string][]types.Release) {
	updateTime := time.Now().UTC()
	db.UpdateLastScrapeTime(updateTime)

	semaphore := make(chan struct{}, 10)

	var wg sync.WaitGroup

	fmt.Println("Uploading releases...")

	for releaseType, releasesOfType := range releases {
		for _, release := range releasesOfType {
			semaphore <- struct{}{}
			wg.Add(1)
			go func(releaseType string, release types.Release) {
				defer func() {
					wg.Done()
					<-semaphore
				}()
				releaseId, err := db.AddOrUpdateRelease(releaseType, release, updateTime)
				if err == nil {
					AddOrUpdateArtists(releaseId, release)
					db.AddOrUpdateProducers(releaseId, release)
					db.AddOrUpdateGenres(releaseId, release)
				}
			}(releaseType, release)
		}
	}

	wg.Wait()

	db.PurgeReleases(updateTime)
}

func AddOrUpdateArtists(releaseId int64, release types.Release) {

	artitstsPopularity := types.PopularityAverage{}
	featuresPopularity := types.PopularityAverage{}

	for _, artist := range release.Artists {
		spotifyArtist, _ := SpotifySearch(artist)

		artistId, popularity, err := db.UploadArtist(artist, spotifyArtist)
		if err == nil {
			if popularity != -1 {
				artitstsPopularity.AddValue(popularity)
			}
			db.UploadReleaseArtists(releaseId, artistId, "main")
		}
	}

	for _, feature := range release.Featurings {
		spotifyArtist, _ := SpotifySearch(feature)

		artistId, popularity, err := db.UploadArtist(feature, spotifyArtist)
		if err == nil {
			if popularity != -1 {
				featuresPopularity.AddValue(popularity)
			}
			db.UploadReleaseArtists(releaseId, artistId, "feature")
		}
	}

	artistsAverage := artitstsPopularity.GetAverage()
	featuresAverage := featuresPopularity.GetAverage()
	var trending_score float64

	if artistsAverage > 0 && featuresAverage > 0 {
		trending_score = (artistsAverage * 0.75) + (featuresAverage * 0.25)
	} else if artistsAverage > 0 {
		trending_score = artistsAverage
	}

	if trending_score != 0.0 {
		db.UpdateReleaseTrendingScore(releaseId, trending_score)
	}
}

func ScrapeReleases() {
	now := time.Now()
	startDate := now.AddDate(0, 0, -14)
	endDate := now.AddDate(0, 3, 0)

	releaseTypes := []string{"lp", "ep", "single", "mixtape", "reissue"}
	allReleases := make(map[string][]types.Release)

	for _, releaseType := range releaseTypes {
		allReleases[releaseType] = []types.Release{}
	}

	var wg sync.WaitGroup
	wg.Add(len(releaseTypes))

	fmt.Println("Fetching releases...")

	for i, releaseType := range releaseTypes {
		requestDate := startDate.Format(requestDateFormat)
		limitRequestDate := endDate.Format(requestDateFormat)

		go func(i int, releaseType string) {
			defer wg.Done()

			for requestDate <= limitRequestDate {
				releases := allReleases[releaseType]
				getReleases(requestDate, startDate, endDate, 0, &releases, releaseType)
				allReleases[releaseType] = releases

				t, _ := time.Parse(requestDateFormat, requestDate)
				t = t.AddDate(0, 1, 0)
				requestDate = t.Format(requestDateFormat)
			}

		}(i, releaseType)
	}

	wg.Wait()

	for i, releaseType := range releaseTypes {
		fmt.Print(releaseType, " count: ", len(allReleases[releaseType]))
		if i != len(releaseTypes)-1 {
			fmt.Print(" / ")
		} else {
			fmt.Print("\n")
		}
	}

	Upload(allReleases)
}

func getReleases(requestDate string, startDate time.Time, endDate time.Time, start int, allReleases *[]types.Release, releaseType string) {
	year := requestDate[:4]

	data := map[string]string{
		"type":      "albumMonth",
		"sort":      "release",
		"albumType": releaseType,
		"start":     fmt.Sprintf("%v", start),
		"date":      requestDate,
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
		date := e.ChildText("div.date") + " " + year
		parsedDate, parseErr := time.Parse(parseDateFormat, date)

		link := e.DOM.Find("div.albumTitle").Parent().AttrOr("href", "")
		aoty_id, idErr := extractReleaseID(link)

		cover := e.ChildAttr("img.lazyload", "data-src")
		title := e.ChildText("div.albumTitle")

		release := types.Release{
			AOTY_Id:    aoty_id,
			Artists:    []string{},
			Featurings: []string{},
			Title:      title,
			Date:       parsedDate,
			Cover:      cover,
			Genres:     []string{},
			Producers:  []string{},
		}

		if parseErr == nil && idErr == nil && !parsedDate.Before(startDate) && !parsedDate.After(endDate) {
			wg.Add(1)
			go func() {
				defer wg.Done()
				getDetails(link, &release)
				*allReleases = append(*allReleases, release)
			}()
		}

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

func extractReleaseID(url string) (string, error) {
	start := strings.Index(url, "/album/") + 7
	if start == -1 {
		return "", errors.New("invalid URL")
	}
	end := strings.Index(url[start:], "-") + start
	if end == -1 {
		return "", errors.New("invalid URL")
	}
	return url[start:end], nil
}

func getDetails(link string, details *types.Release) {
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

	re1 := regexp.MustCompile(`(?i)(?:(\[featuring|\(featuring|\s+featuring)\.?\s|(\[feat|\(feat|\s+feat)\.?\s|(\[ft|\(ft|\s+ft)\.?\s)`)
	matchIndexes1 := re1.FindStringSubmatchIndex(songTitle)

	re2 := regexp.MustCompile(`(?i)(?:(\[with|\(with|\(w/|\[w/|\s+w/\s)\.?)`)
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
	re := regexp.MustCompile(`(?i)(?:(\[prod. by|\(prod. by|\s+prod. by|\[prod by|\(prod by|\s+prod by|\[prod.|\(prod.|\s+prod.|\[prod|\(prod|\s+prod\s))`)
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
	parts := SplitString(titleSection, separators)
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" && part != " " && part != "." {
			part = cleanContributorName(part)
			*contributors = append(*contributors, part)
		}
	}
}

func SplitString(str string, separators []string) []string {
	if len(separators) == 0 {
		return []string{str}
	}
	sep := separators[0]
	parts := strings.Split(str, sep)
	result := []string{}
	for _, part := range parts {
		result = append(result, SplitString(part, separators[1:])...)
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
