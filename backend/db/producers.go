package db

import (
	"database/sql"
	"fmt"
	"main/types"
	"strings"
)

func AddOrUpdateProducers(releaseId int64, release types.Release) {
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
				fmt.Println("Error inserting release producer join")
			}
		}
	}
}
