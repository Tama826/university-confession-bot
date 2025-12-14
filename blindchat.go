package main

import (
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// ---------------- BLIND CHAT ----------------

func startBlindChat(uid int64) {
	if _, ok := Pairs[uid]; ok {
		bot.Send(tgbotapi.NewMessage(uid, "⚠️ You are already in a chat."))
		return
	}

	if WaitingUser == 0 {
		WaitingUser = uid
		WaitingSince = time.Now().Unix()
		bot.Send(tgbotapi.NewMessage(uid, "⏳ Waiting for a partner..."))
		return
	}

	if WaitingUser == uid {
		return
	}

	Pairs[uid] = WaitingUser
	Pairs[WaitingUser] = uid

	bot.Send(tgbotapi.NewMessage(uid, "🔗 Connected anonymously. Say hi!"))
	bot.Send(tgbotapi.NewMessage(WaitingUser, "🔗 Connected anonymously. Say hi!"))

	WaitingUser = 0
	WaitingSince = 0
}

func endBlindChat(uid int64) {
	partner, ok := Pairs[uid]
	if !ok {
		bot.Send(tgbotapi.NewMessage(uid, "❌ You are not in a blind chat."))
		return
	}

	delete(Pairs, uid)
	delete(Pairs, partner)

	bot.Send(tgbotapi.NewMessage(uid, "❌ Chat ended."))
	bot.Send(tgbotapi.NewMessage(partner, "❌ Chat ended."))
}

func reportPartner(uid int64) {
	partner, ok := Pairs[uid]
	if !ok {
		bot.Send(tgbotapi.NewMessage(uid, "❌ No partner to report."))
		return
	}

	Reports[partner]++
	bot.Send(tgbotapi.NewMessage(uid, "🚨 Report submitted."))
	bot.Send(tgbotapi.NewMessage(AdminGroupID, "🚨 Blind chat abuse reported."))

	if Reports[partner] >= 3 {
		db.Exec("UPDATE users SET banned=1 WHERE user_id=?", partner)
		bot.Send(tgbotapi.NewMessage(partner, "⛔ You have been banned due to reports."))
		endBlindChat(partner)
	}
}
