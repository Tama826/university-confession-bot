package main

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func initDatabase() {
	var err error
	db, err = sql.Open("sqlite3", "confess.db")
	if err != nil {
		log.Fatal(err)
	}

	createTables()
}

func createTables() {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			user_id INTEGER PRIMARY KEY,
			role TEXT DEFAULT 'user',
			banned INTEGER DEFAULT 0,
			reports INTEGER DEFAULT 0
		);`,
		`CREATE TABLE IF NOT EXISTS confessions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			text TEXT,
			status TEXT DEFAULT 'pending',
			created INTEGER,
			scheduled INTEGER DEFAULT 0,
			channel_msg_id INTEGER DEFAULT 0,
			upvotes INTEGER DEFAULT 0,
			downvotes INTEGER DEFAULT 0
		);`,
		`CREATE TABLE IF NOT EXISTS votes (
			confession_id INTEGER,
			user_id INTEGER,
			vote INTEGER,
			PRIMARY KEY(confession_id, user_id)
		);`,
		`CREATE TABLE IF NOT EXISTS keywords (
			word TEXT PRIMARY KEY,
			autoban INTEGER DEFAULT 0
		);`,
		`CREATE TABLE IF NOT EXISTS logs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			action TEXT,
			timestamp INTEGER
		);`,
	}

	for _, q := range queries {
		if _, err := db.Exec(q); err != nil {
			log.Fatal(err)
		}
	}
}

// Utility to get current timestamp
func nowUnix() int64 {
	return time.Now().Unix()
}
