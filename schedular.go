package main

import (
	"time"
)

// ---------------- SCHEDULER ----------------

func startScheduler() {
	go func() {
		for {
			rows, _ := db.Query(`
				SELECT id FROM confessions 
				WHERE status='pending' AND scheduled>0 AND scheduled<=?
			`, nowUnix())

			for rows.Next() {
				var id int64
				rows.Scan(&id)
				publishConfession(id)
			}

			time.Sleep(ScheduleCheckInterval)
		}
	}()
}

func scheduleConfession(confID int64, t int64) {
	db.Exec("UPDATE confessions SET scheduled=? WHERE id=?", t, confID)
}
