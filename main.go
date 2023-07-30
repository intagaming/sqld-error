package main

import (
	"database/sql"
	_ "github.com/libsql/libsql-client-go/libsql"
	"log"
	"os"
	"time"
)

func main() {
	var dbUrl = os.Getenv("DATABASE_URL")
	db, err := sql.Open("libsql", dbUrl)
	if err != nil {
		log.Fatalf("failed to open db %s: %v", dbUrl, err)
	}

	statements := []string{
		`CREATE TABLE IF NOT EXISTS user (
			sub TEXT PRIMARY KEY
		)`,
		`CREATE TABLE IF NOT EXISTS stream (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id TEXT NOT NULL,
			url TEXT NOT NULL,
			platform TEXT NOT NULL,
			created_at INTEGER NOT NULL,
			gone_online INTEGER DEFAULT 0 NOT NULL,
			terminated_at INTEGER,
			scheduled_end_at INTEGER,
			error TEXT,
			close_reason TEXT,
			password TEXT NOT NULL,

			FOREIGN KEY(user_id) REFERENCES user(sub)
		)`,
	}
	for _, s := range statements {
		_, err := db.Exec(s)
		if err != nil {
			log.Fatalf("failed executing initiate db statement '%s': %v", s, err)
		}
	}

	log.Printf("inserting")
	_, err = db.Exec("INSERT OR IGNORE INTO user (sub) VALUES (?)", "test")
	if err != nil {
		panic(err)
	}

	go func(innerDb *sql.DB) {
		for {
			go func() {
				log.Printf("querying")
				sub := "test"
				rows, err := innerDb.Query("SELECT sub FROM user WHERE sub = ?", sub)
				if err != nil {
					log.Fatalf("failed querying: %v", err)
					panic(err)
				}
				defer func() {
                    _ = rows.Close()
					// if err != nil {
					// 	panic(err)
					// }
					// err = rows.Err()
					// if err != nil {
					// 	panic(err)
					// }
				}()
				if rows.Next() {
					// var sub string
					// rows.Scan(&sub)
					// log.Printf("found %s", sub)
				}
				// err = rows.Close()
				// if err != nil {
				// 	panic(err)
				// }
				// err = rows.Err()
				// if err != nil {
				// 	panic(err)
				// }
			}()
			// go func() {
			// 	log.Printf("inserting")
			// 	_, err := innerDb.Exec("INSERT OR IGNORE INTO user (sub) VALUES (?)", "test")
			// 	if err != nil {
			// 		panic(err)
			// 	}
			// }()
			// go func() {
			// 	log.Printf("updating")
			// 	_, err := innerDb.Exec("UPDATE user SET sub = ? WHERE sub = ?", "test", "test")
			// 	if err != nil {
			// 		panic(err)
			// 	}
			// }()
			// time.Sleep(200 * time.Millisecond)
			time.Sleep(15 * time.Second)
		}
	}(db)
	for {
	}
}
