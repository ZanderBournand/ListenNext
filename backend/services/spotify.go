package services

import (
	"encoding/json"
	"fmt"
	"log"
	"main/models"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/caffix/cloudflare-roundtripper/cfrt"
	"github.com/gocolly/colly"
)

var (
	scrapingClientIDs = []string{
		os.Getenv("SPOTIFY_API_SCRAPING_CLIENT_IDS"),
		os.Getenv("SPOTIFY_API_SCRAPING_CLIENT_IDS"),
		os.Getenv("SPOTIFY_API_SCRAPING_CLIENT_IDS"),
	}
	scrapingClientSecrets = []string{
		os.Getenv("SPOTIFY_API_GENERAL_CLIENT_SECRET"),
		os.Getenv("SPOTIFY_API_GENERAL_CLIENT_SECRET"),
		os.Getenv("SPOTIFY_API_GENERAL_CLIENT_SECRET"),
	}
	generalClientID     = os.Getenv("SPOTIFY_API_GENERAL_CLIENT_ID")
	generalClientSecret = os.Getenv("SPOTIFY_API_GENERAL_CLIENT_SECRET")
)

var (
	scrapingTokens []string
	generalToken   string
)

func SpotifyScrapeTokens() {
	tokens := make([]string, 3)

	for i := 0; i < 3; i++ {
		req, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", strings.NewReader("grant_type=client_credentials"))
		if err != nil {
			panic(err)
		}

		req.SetBasicAuth(scrapingClientIDs[i], scrapingClientSecrets[i])
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		var data struct {
			AccessToken string `json:"access_token"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			panic(err)
		}

		tokens[i] = data.AccessToken
	}

	scrapingTokens = tokens
}

func SpotifyGeneralToken() {
	req, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", strings.NewReader("grant_type=client_credentials"))
	if err != nil {
		panic(err)
	}

	req.SetBasicAuth(generalClientID, generalClientSecret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var data struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		panic(err)
	}

	generalToken = data.AccessToken
}

func SpotifyUserTops(client *http.Client) ([]models.SpotifyArtist, []models.SpotifyTrack, error) {
	limitQuery := "10"

	var spotifyArtists []models.SpotifyArtist
	var spotifyTracks []models.SpotifyTrack

	endpoint := fmt.Sprintf("https://api.spotify.com/v1/me/top/artists?limit=%s&time_range=short_term", limitQuery)

	resp, err := client.Get(endpoint)
	if err != nil {
		return nil, nil, fmt.Errorf("error making user tops request: %v", err)
	}
	defer resp.Body.Close()

	var data map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, nil, fmt.Errorf("error decoding response: %v", err)
	}

	artists := data["items"].([]interface{})

	for _, artist := range artists {
		var spotifyArtist models.SpotifyArtist

		spotifyArtist.Name = artist.(map[string]interface{})["name"].(string)
		spotifyArtist.ID = artist.(map[string]interface{})["id"].(string)
		genres := make([]string, 0)
		for _, genre := range artist.(map[string]interface{})["genres"].([]interface{}) {
			genres = append(genres, genre.(string))
		}
		spotifyArtist.Genres = genres

		spotifyArtists = append(spotifyArtists, spotifyArtist)
	}

	endpoint = fmt.Sprintf("https://api.spotify.com/v1/me/top/tracks?limit=%s&time_range=short_term", limitQuery)
	resp, err = client.Get(endpoint)
	if err != nil {
		return nil, nil, fmt.Errorf("error making user tops request: %v", err)
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, nil, fmt.Errorf("error decoding response: %v", err)
	}

	items := data["items"].([]interface{})

	var wg sync.WaitGroup

	for _, item := range items {
		wg.Add(1)
		go func(item interface{}) {
			defer wg.Done()

			var spotifyTrack models.SpotifyTrack
			spotifyTrack.ID = item.(map[string]interface{})["id"].(string)

			artists := item.(map[string]interface{})["artists"].([]interface{})

			for _, artist := range artists {

				var spotifyArtist models.SpotifyArtist

				spotifyArtist.Name = artist.(map[string]interface{})["name"].(string)
				spotifyArtist.ID = artist.(map[string]interface{})["id"].(string)

				endpoint := fmt.Sprintf("https://api.spotify.com/v1/artists/%s", spotifyArtist.ID)

				rand.Seed(time.Now().UnixNano())
				tokenIndex := rand.Intn(len(scrapingTokens))

				req, err := http.NewRequest("GET", endpoint, nil)
				if err != nil {
					log.Printf("error making artist request: %v", err)
					return
				}
				req.Header.Set("Authorization", "Bearer "+scrapingTokens[tokenIndex])

				client := &http.Client{}
				resp, err := client.Do(req)
				if err != nil || resp.StatusCode != 200 {
					log.Printf("error making artist request: %v", err)
					return
				}
				defer resp.Body.Close()

				var artistData map[string]interface{}
				err = json.NewDecoder(resp.Body).Decode(&artistData)
				if err != nil {
					log.Printf("error decoding artist response: %v", err)
					return
				}

				genres := artistData["genres"].([]interface{})
				for _, genre := range genres {
					spotifyArtist.Genres = append(spotifyArtist.Genres, genre.(string))
				}

				spotifyArtists = append(spotifyArtists, spotifyArtist)
			}

			spotifyTracks = append(spotifyTracks, spotifyTrack)
		}(item)
	}

	wg.Wait()

	return spotifyArtists, spotifyTracks, nil
}

func SpotifyRelatedArtists(artists []models.SpotifyArtist) ([]models.SpotifyArtist, error) {
	var wg sync.WaitGroup

	for _, artist := range artists {
		wg.Add(1)
		go func(artist models.SpotifyArtist) {
			defer wg.Done()
			endpoint := fmt.Sprintf("https://api.spotify.com/v1/artists/%s/related-artists", artist.ID)

			rand.Seed(time.Now().UnixNano())
			tokenIndex := rand.Intn(len(scrapingTokens))

			req, err := http.NewRequest("GET", endpoint, nil)
			if err != nil {
				panic(err)
			}
			req.Header.Set("Authorization", "Bearer "+scrapingTokens[tokenIndex])

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil || resp.StatusCode != 200 {
				panic(err)
			}
			defer resp.Body.Close()

			var data map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&data)
			if err != nil {
				panic(err)
			}

			relatedArtistsData := data["artists"].([]interface{})

			for _, relatedArtistData := range relatedArtistsData {
				relatedArtist := models.SpotifyArtist{
					Name: relatedArtistData.(map[string]interface{})["name"].(string),
					ID:   relatedArtistData.(map[string]interface{})["id"].(string),
				}
				genres := make([]string, 0)
				for _, genre := range relatedArtistData.(map[string]interface{})["genres"].([]interface{}) {
					genres = append(genres, genre.(string))
				}
				relatedArtist.Genres = genres

				artists = append(artists, relatedArtist)
			}
		}(artist)
	}

	wg.Wait()

	return artists, nil
}

func SpotifyRecommendations(artists []models.SpotifyArtist, tracks []models.SpotifyTrack) ([]models.SpotifyArtist, error) {
	posArtists := 0
	posTracks := 0

	artistsCount := len(artists)
	tracksCount := len(tracks)

	var wg sync.WaitGroup

	for posArtists < artistsCount || posTracks < tracksCount {

		var artistIDs []string
		if posArtists < artistsCount {
			// numArtists := rand.Intn(5) + 1
			// if posTracks == tracksCount {
			// 	numArtists = 5
			// }
			numArtists := 5
			for i := 0; i < numArtists && posArtists < artistsCount; i++ {
				artistIDs = append(artistIDs, artists[posArtists].ID)
				posArtists++
			}
		}

		var trackIDs []string
		if posTracks < tracksCount {
			numTracks := 5 - len(artistIDs)
			for i := 0; i < numTracks && posTracks < tracksCount; i++ {
				trackIDs = append(trackIDs, tracks[posTracks].ID)
				posTracks++
			}
		}

		wg.Add(1)

		go func() {
			defer wg.Done()

			recommendations, err := Recommendations(artistIDs, nil, trackIDs)
			if err != nil {
				log.Printf("error getting recommendations: %v", err)
				return
			}

			for _, artist := range recommendations {
				found := false
				for _, existingArtist := range artists {
					if existingArtist.ID == artist.ID {
						found = true
						break
					}
				}
				if !found {
					artists = append(artists, artist)
				}
			}
		}()
	}

	wg.Wait()

	return artists, nil
}

func Recommendations(artistIds []string, genres []string, trackIds []string) ([]models.SpotifyArtist, error) {
	limitQuery := "10"

	params := url.Values{}
	if len(artistIds) > 0 {
		params.Set("seed_artists", strings.Join(artistIds, ","))
	}
	if len(genres) > 0 {
		params.Set("seed_genres", strings.Join(genres, ","))
	}
	if len(trackIds) > 0 {
		params.Set("seed_tracks", strings.Join(trackIds, ","))
	}

	endpoint := "https://api.spotify.com/v1/recommendations?limit=" + limitQuery + "&" + params.Encode()

	rand.Seed(time.Now().UnixNano())
	tokenIndex := rand.Intn(len(scrapingTokens))

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+scrapingTokens[tokenIndex])

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var artistIDs []string
	var recommendations struct {
		Tracks []struct {
			Artists []struct {
				ID string `json:"id"`
			} `json:"artists"`
		} `json:"tracks"`
	}

	err = json.NewDecoder(res.Body).Decode(&recommendations)
	if err != nil {
		return nil, err
	}

	for _, track := range recommendations.Tracks {
		for _, artist := range track.Artists {
			artistIDs = append(artistIDs, artist.ID)
		}
	}

	var releaseArtists []models.SpotifyArtist
	var wg sync.WaitGroup

	for _, artistID := range artistIDs {
		wg.Add(1)
		go func(artistID string) {
			defer wg.Done()

			var spotifyArtist models.SpotifyArtist
			spotifyArtist.ID = artistID

			endpoint := fmt.Sprintf("https://api.spotify.com/v1/artists/%s", spotifyArtist.ID)

			rand.Seed(time.Now().UnixNano())
			tokenIndex := rand.Intn(len(scrapingTokens))

			req, err := http.NewRequest("GET", endpoint, nil)
			if err != nil {
				panic(err)
			}
			req.Header.Set("Authorization", "Bearer "+scrapingTokens[tokenIndex])

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				panic(err)
			}
			defer resp.Body.Close()

			var artistData map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&artistData)
			if err != nil {
				panic(err)
			}

			spotifyArtist.Name = artistData["name"].(string)

			genres := artistData["genres"].([]interface{})
			for _, genre := range genres {
				spotifyArtist.Genres = append(spotifyArtist.Genres, genre.(string))
			}

			releaseArtists = append(releaseArtists, spotifyArtist)
		}(artistID)
	}

	wg.Wait()

	return releaseArtists, nil
}

func TopGenres(artists []models.SpotifyArtist) []string {
	totalGenres := 20

	genreCount := make(map[string]int)
	for _, artist := range artists {
		for _, genre := range artist.Genres {
			genreCount[genre]++
		}
	}

	sortedGenres := make([]string, 0, len(genreCount))
	for genre := range genreCount {
		sortedGenres = append(sortedGenres, genre)
	}
	sort.Slice(sortedGenres, func(i, j int) bool {
		return genreCount[sortedGenres[i]] > genreCount[sortedGenres[j]]
	})

	topGenres := make([]string, 0, totalGenres)
	for i := 0; i < totalGenres && i < len(sortedGenres); i++ {
		topGenres = append(topGenres, sortedGenres[i])
	}

	return topGenres
}

func SpotifySearch(artist string) (string, string, int, []string, error) {
	compareName := strings.ToLower(strings.ReplaceAll(artist, " ", ""))

	spotifyURL := fmt.Sprintf("https://api.spotify.com/v1/search?type=artist&q=%s", url.QueryEscape(artist))

	rand.Seed(time.Now().UnixNano())
	tokenIndex := rand.Intn(len(scrapingTokens))

	req, err := http.NewRequest("GET", spotifyURL, nil)
	if err != nil {
		return "", "", -1, nil, err
	}
	req.Header.Set("Authorization", "Bearer "+scrapingTokens[tokenIndex])

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != 200 {
		fmt.Println("SPOTIFY ERROR!!!!")
		fmt.Println(resp.Header)
		return "", "", -1, nil, err
	}
	defer resp.Body.Close()

	var data map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return "", "", -1, nil, err
	}

	artists := data["artists"].(map[string]interface{})["items"].([]interface{})
	for _, artist := range artists {
		name := strings.ToLower(strings.ReplaceAll(artist.(map[string]interface{})["name"].(string), " ", ""))
		if name == compareName {
			genres := make([]string, 0)
			for _, genre := range artist.(map[string]interface{})["genres"].([]interface{}) {
				genres = append(genres, genre.(string))
			}
			return artist.(map[string]interface{})["name"].(string), artist.(map[string]interface{})["id"].(string), int(artist.(map[string]interface{})["popularity"].(float64)), genres, nil
		}
	}

	return "", "", -1, nil, fmt.Errorf("no matching artist found")
}

func SpotifyScrapedSearch(artist string) (string, string, int, []string, error) {
	compareName := strings.ToLower(strings.ReplaceAll(artist, " ", ""))
	found := false

	var name string
	var spotifyId string
	var popularity int
	var genres []string

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

	c.OnHTML("div.song-details.search-song-details", func(e *colly.HTMLElement) {
		spotifyName := e.ChildText("h1.song-title u")
		spotifyCompareName := strings.ToLower(strings.ReplaceAll(spotifyName, " ", ""))

		if spotifyCompareName == compareName && !found {
			found = true

			name = spotifyName

			href := e.ChildAttr("a:nth-of-type(1)", "href")
			parts := strings.Split(href, "/")
			spotifyId = parts[len(parts)-1]

			c.Visit("https://musicstax.com/" + href)
		}
	})

	c.OnHTML("div.song-details-right", func(e *colly.HTMLElement) {
		allGenres := e.ChildText("[data-cy='artist-genres']")
		separators := []string{", ", " & "}
		genres = SplitString(allGenres, separators)

		popularityStr := strings.TrimSpace(strings.Split(e.ChildText(`[data-cy="artist-followers"]`), "//")[1])
		popularityParsed, err := strconv.Atoi(strings.TrimSuffix(popularityStr, "% popularity"))
		if err != nil {
			log.Println("Error parsing popularity:", err)
			return
		}
		popularity = popularityParsed
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("ERROR:", r.StatusCode)
		fmt.Println(r.Headers)
	})

	c.Visit("https://musicstax.com/search?q=" + artist + "&view=artists")
	c.Wait()

	if found {
		return name, spotifyId, popularity, genres, nil
	} else {
		return "", "", -1, nil, fmt.Errorf("no matching artist found")
	}
}
