package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/lib/pq"
)

func Upload(releases map[string][]Release, spotifyAuthToken string, db *sql.DB) {
	updateTime := time.Now()

	semaphore := make(chan struct{}, 50)

	var wg sync.WaitGroup

	for releaseType, releasesOfType := range releases {
		fmt.Println("Uploading " + releaseType + "s...")
		for _, release := range releasesOfType {
			semaphore <- struct{}{}
			wg.Add(1)
			go func(releaseType string, release Release) {
				defer func() {
					wg.Done()
					<-semaphore
				}()
				releaseId, err := AddOrUpdateRelease(releaseType, release, db, updateTime)
				if err == nil {
					AddOrUpdateArtists(releaseId, release, spotifyAuthToken, db)
					AddOrUpdateProducers(releaseId, release, db)
					AddOrUpdateGenres(releaseId, release, db)
				}
			}(releaseType, release)
		}
	}

	wg.Wait()

	fmt.Println("----------------------------")
	fmt.Println("----------------------------")
	_, err := db.Exec("DELETE FROM Releases WHERE updated!=$1", updateTime)
	if err == nil {
		fmt.Println("Old releases purged!")
	}
}

func AddOrUpdateRelease(releaseType string, release Release, db *sql.DB, updateTime time.Time) (int64, error) {
	var id int64
	err := db.QueryRow("SELECT id FROM Releases WHERE aoty_id=$1", release.AOTY_Id).Scan(&id)
	if err == nil {
		_, err = db.Exec("UPDATE Releases SET title=$1, date=$2, cover=$3, tracklist=$4, type=$5, updated=$6 WHERE id=$7",
			release.Title, release.Date, release.Cover, pq.Array(release.Tracklist), releaseType, updateTime, id)
		if err == nil {
			return id, nil
		}
		fmt.Println("error updating release")
		return -1, fmt.Errorf("error updating release")
	} else if err == sql.ErrNoRows {
		var newID int64
		err = db.QueryRow("INSERT INTO Releases(title, date, cover, tracklist, type, updated, aoty_id) VALUES($1, $2, $3, $4, $5, $6, $7) RETURNING id",
			release.Title, release.Date, release.Cover, pq.Array(release.Tracklist), releaseType, updateTime, release.AOTY_Id).Scan(&newID)
		if err == nil {
			return newID, nil
		}
		fmt.Println("error addding release")
		return -1, fmt.Errorf("error addding release")
	} else {
		fmt.Println(err)
		fmt.Println("error querying release")
		return -1, fmt.Errorf("error querying release")
	}

}

func AddOrUpdateArtists(releaseId int64, release Release, spotifyAuthToken string, db *sql.DB) {

	for _, artist := range release.Artists {
		artistId, err := uploadArtist(artist, spotifyAuthToken, db)
		if err == nil {
			_, err := db.Exec("INSERT INTO releases_artists (release_id, artist_id, relationship) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING", releaseId, artistId, "main")
			if err != nil {
				fmt.Println(release)
				fmt.Println("Error inserting release main artist join")
			}
		}
	}

	for _, featuredArtists := range release.Featurings {
		artistId, err := uploadArtist(featuredArtists, spotifyAuthToken, db)
		if err == nil {
			_, err := db.Exec("INSERT INTO releases_artists (release_id, artist_id, relationship) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING", releaseId, artistId, "feature")
			if err != nil {
				fmt.Println(release)
				fmt.Println("Error inserting release feature artist join")
			}
		}
	}
}

func uploadArtist(artist string, spotifyAuthToken string, db *sql.DB) (int64, error) {
	compareName := strings.ToLower(strings.ReplaceAll(artist, " ", ""))

	spotifyName, spotifyId, popularity, genres, err := spotifySearch(artist, spotifyAuthToken)
	if err == nil {
		var id int64
		err := db.QueryRow("SELECT id FROM Artists WHERE name=$1 AND spotify_id=$2", spotifyName, spotifyId).Scan(&id)
		if err == nil {
			_, err = db.Exec("UPDATE Artists SET popularity=$1 WHERE name=$2 AND spotify_id=$3",
				popularity, spotifyName, spotifyId)
			if err == nil {
				AddOrUpdateArtistGenres(id, genres, db)
				return id, nil
			}
			return -1, fmt.Errorf("error updating artist w/ spotify")
		} else if err == sql.ErrNoRows {
			var newID int64
			err = db.QueryRow("INSERT INTO Artists(name, spotify_id, popularity, compare_name) VALUES($1, $2, $3, $4) RETURNING id",
				spotifyName, spotifyId, popularity, compareName).Scan(&newID)
			if err == nil {
				AddOrUpdateArtistGenres(newID, genres, db)
				return newID, nil
			}
			return -1, fmt.Errorf("error inserting artist w/ spotify")
		} else {
			return -1, fmt.Errorf("error querying artist w/ spotify")
		}
	} else {
		compareName := strings.ToLower(strings.ReplaceAll(artist, " ", ""))

		var id int64
		err := db.QueryRow("SELECT id FROM Artists WHERE compare_name=$1", compareName).Scan(&id)
		if err == nil {
			return id, nil
		} else if err == sql.ErrNoRows {
			var newID int64
			err = db.QueryRow("INSERT INTO Artists(name, compare_name) VALUES($1, $2) RETURNING id",
				artist, compareName).Scan(&newID)
			if err == nil {
				return newID, nil
			}
			return -1, fmt.Errorf("error inserting artist")
		} else {
			return -1, fmt.Errorf("error querying artist")
		}
	}
}

func spotifySearch(artist string, spotifyAuthToken string) (string, string, int, []string, error) {
	compareName := strings.ToLower(strings.ReplaceAll(artist, " ", ""))

	spotifyURL := fmt.Sprintf("https://api.spotify.com/v1/search?type=artist&q=%s", url.QueryEscape(artist))

	req, err := http.NewRequest("GET", spotifyURL, nil)
	if err != nil {
		return "", "", -1, nil, err
	}
	req.Header.Set("Authorization", "Bearer "+spotifyAuthToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("SPOTIFY ERROR!!!!")
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

func AddOrUpdateProducers(releaseId int64, release Release, db *sql.DB) {

	for _, producer := range release.Producers {
		compareName := strings.ToLower(strings.ReplaceAll(producer, " ", ""))

		var id int64
		err := db.QueryRow("SELECT id FROM Producers WHERE compare_name=$1", compareName).Scan(&id)
		if err == sql.ErrNoRows {
			_ = db.QueryRow("INSERT INTO Producers(name, compare_name) VALUES($1, $2) RETURNING id",
				producer, compareName).Scan(&id)
		}

		if id != 0 {
			_, err := db.Exec("INSERT INTO releases_producers (release_id, producer_id) VALUES ($1, $2) ON CONFLICT DO NOTHING", releaseId, id)
			if err != nil {
				fmt.Println(release)
				fmt.Println("Error inserting release producer join")
			}
		}
	}

}

func AddOrUpdateGenres(releaseId int64, release Release, db *sql.DB) {

	for _, genre := range release.Genres {
		compareType := strings.ToLower(strings.ReplaceAll(genre, " ", ""))

		var id int64
		err := db.QueryRow("SELECT id FROM Genres WHERE compare_type=$1", compareType).Scan(&id)
		if err == sql.ErrNoRows {
			_ = db.QueryRow("INSERT INTO Genres(type, compare_type) VALUES($1, $2) RETURNING id",
				genre, compareType).Scan(&id)
		}

		if id != 0 {
			_, err := db.Exec("INSERT INTO releases_genres (release_id, genre_id) VALUES ($1, $2) ON CONFLICT DO NOTHING", releaseId, id)
			if err != nil {
				fmt.Println(release)
				fmt.Println("Error inserting release genre join")
			}
		}
	}

}

func AddOrUpdateArtistGenres(artistId int64, genres []string, db *sql.DB) {

	for _, genre := range genres {
		compareType := strings.ToLower(strings.ReplaceAll(genre, " ", ""))

		var id int64
		err := db.QueryRow("SELECT id FROM Genres WHERE compare_type=$1", compareType).Scan(&id)
		if err == sql.ErrNoRows {
			_ = db.QueryRow("INSERT INTO Genres(type, compare_type) VALUES($1, $2) RETURNING id",
				genre, compareType).Scan(&id)
		}

		if id != 0 {
			_, err := db.Exec("INSERT INTO artists_genres (artist_id, genre_id) VALUES ($1, $2) ON CONFLICT DO NOTHING", artistId, id)
			if err != nil {
				fmt.Println(artistId)
				fmt.Println("Error inserting artist genre join")
			}
		}
	}

}
