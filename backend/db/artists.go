package db

import (
	"database/sql"
	"fmt"
	"main/types"
	"strings"
)

func GetMatchingArtists(artists []types.SpotifyArtist, genres []string) []int {
	var ids []int

	var allIDs []string
	var allGenres []string

	for _, artist := range artists {
		allIDs = append(allIDs, "'"+artist.ID+"'")
	}

	for _, genre := range genres {
		compareType := strings.ToLower(strings.ReplaceAll(genre, " ", ""))
		allGenres = append(allGenres, "'"+compareType+"'")
	}

	query := `
	    SELECT id
	    FROM Artists
	    WHERE spotify_id IN (` + strings.Join(allIDs, ",") + `)
	        OR EXISTS (
	            SELECT 1
	            FROM Artists_Genres
	            JOIN Genres ON Artists_Genres.genre_id = Genres.id
	            WHERE Artists_Genres.artist_id = Artists.id
	                AND Genres.compare_type IN (` + strings.Join(allGenres, ",") + `)
	        )
	`

	rows, err := db.Query(query)
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			fmt.Println(err)
		}
		ids = append(ids, id)
	}
	if err := rows.Err(); err != nil {
		fmt.Println(err)
	}

	return ids
}

func UploadArtist(artist string, spotifyArtist *types.SpotifyArtist) (int64, int, error) {
	compareName := strings.ToLower(strings.ReplaceAll(artist, " ", ""))

	if spotifyArtist != nil {
		var id int64
		err := db.QueryRow("SELECT id FROM Artists WHERE name=$1 AND spotify_id=$2", spotifyArtist.Name, spotifyArtist.ID).Scan(&id)
		if err == nil {
			_, err = db.Exec("UPDATE Artists SET popularity=$1 WHERE name=$2 AND spotify_id=$3",
				spotifyArtist.Popularity, spotifyArtist.Name, spotifyArtist.ID)
			if err == nil {
				AddOrUpdateArtistGenres(id, spotifyArtist.Genres)
				return id, spotifyArtist.Popularity, nil
			}
			return -1, -1, fmt.Errorf("error updating artist w/ spotify")
		} else if err == sql.ErrNoRows {
			var newID int64
			err = db.QueryRow("INSERT INTO Artists(name, spotify_id, popularity, compare_name) VALUES($1, $2, $3, $4) RETURNING id",
				spotifyArtist.Name, spotifyArtist.ID, spotifyArtist.Popularity, compareName).Scan(&newID)
			if err == nil {
				AddOrUpdateArtistGenres(newID, spotifyArtist.Genres)
				return newID, spotifyArtist.Popularity, nil
			}
			return -1, -1, fmt.Errorf("error inserting artist w/ spotify")
		} else {
			return -1, -1, fmt.Errorf("error querying artist w/ spotify")
		}
	} else {
		compareName := strings.ToLower(strings.ReplaceAll(artist, " ", ""))

		var id int64
		err := db.QueryRow("SELECT id FROM Artists WHERE compare_name=$1", compareName).Scan(&id)
		if err == nil {
			return id, -1, nil
		} else if err == sql.ErrNoRows {
			var newID int64
			err = db.QueryRow("INSERT INTO Artists(name, compare_name) VALUES($1, $2) RETURNING id",
				artist, compareName).Scan(&newID)
			if err == nil {
				return newID, -1, nil
			}
			return -1, -1, fmt.Errorf("error inserting artist")
		} else {
			return -1, -1, fmt.Errorf("error querying artist")
		}
	}
}

func AddOrUpdateArtistGenres(artistId int64, genres []string) {

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
