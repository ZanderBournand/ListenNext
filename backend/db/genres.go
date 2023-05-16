package db

import (
	"database/sql"
	"fmt"
	"main/types"
	"strings"
)

func AddOrUpdateGenres(releaseId int64, release types.Release) {
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
