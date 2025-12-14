package main

import (
	"fmt"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// ---------------- MESSAGE HANDLER ----------------

func handleMessage(msg *tgbotapi.Message) {
	uid := msg.From.ID

	// Timeout waiting user for blind chat
	if WaitingUser != 0 && time.Now().Unix()-WaitingSince > 120 {
		bot.Send(tgbotapi.NewMessage(WaitingUser, "❌ No partner found. Try again later."))
		WaitingUser = 0
	}

	// Ignore banned users
	if isBanned(uid) {
		return
	}

	// Commands
	if msg.IsCommand() {
		switch msg.Command() {
		case "start":
			bot.Send(tgbotapi.NewMessage(msg.Chat.ID,
				"Welcome!\n/confess – send confession\n/blind – anonymous chat\n/end – end chat\n/report – report abuse"))

		case "confess":
			bot.Send(tgbotapi.NewMessage(msg.Chat.ID, "Send your confession."))

		case "blind":
			startBlindChat(uid)

		case "end":
			endBlindChat(uid)

		case "report":
			reportPartner(uid)
		}
		return
	}

	// Blind chat relay
	if partner, ok := Pairs[uid]; ok {
		bot.Send(tgbotapi.NewMessage(partner, msg.Text))
		return
	}

	// Private confession
	if msg.Chat.IsPrivate() {
		if !canSend(uid) {
			bot.Send(tgbotapi.NewMessage(msg.Chat.ID, "Limit reached. Try later."))
			return
		}

		// Auto-ban keywords / toxicity
		if containsBannedWord(msg.Text) || toxicityScore(msg.Text) >= 60 {
			db.Exec("UPDATE users SET banned=1 WHERE user_id=?", uid)
			bot.Send(tgbotapi.NewMessage(msg.Chat.ID, "⛔ You have been banned due to policy violation."))
			return
		}

		saveUser(uid)
		sendToAdmin(uid, msg.Text)
	}
}

// ---------------- UTILITY FUNCTIONS ----------------

func canSend(uid int64) bool {
	var last int64
	var lastDate string
	db.QueryRow("SELECT last_confession, created FROM users WHERE user_id=?", uid).Scan(&lastDate, &last)
	today := time.Now().Format("2006-01-02")
	if lastDate == today {
		return false
	}
	if time.Now().Unix()-last < 30 {
		return false
	}
	return true
}

func saveUser(uid int64) {
	now := time.Now()
	db.Exec("INSERT OR REPLACE INTO users(user_id,last_confession,created) VALUES(?,?,?)",
		uid, now.Format("2006-01-02"), now.Unix())
}

func isBanned(uid int64) bool {
	var banned int
	db.QueryRow("SELECT banned FROM users WHERE user_id=?", uid).Scan(&banned)
	return banned == 1
}
