package main

import (
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// ---------------- ADMIN INLINE PANEL ----------------

func sendToAdmin(uid int64, text string) {
	res, _ := db.Exec("INSERT INTO confessions(text,created) VALUES(?,?)", text, nowUnix())
	confID, _ := res.LastInsertId()

	msg := tgbotapi.NewMessage(
		AdminGroupID,
		fmt.Sprintf("📩 New Confession (ID %d)\n\n%s\n\n👤 Sender: %d", confID, text, uid),
	)
	msg.ReplyMarkup = adminKeyboard(confID, uid)
	bot.Send(msg)
}

// Admin keyboard
func adminKeyboard(confID int64, uid int64) *tgbotapi.InlineKeyboardMarkup {
	return &tgbotapi.InlineKeyboardMarkup{
		InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{
			{
				tgbotapi.NewInlineKeyboardButtonData("✅ Approve", fmt.Sprintf("approve:%d:%d", confID, uid)),
				tgbotapi.NewInlineKeyboardButtonData("❌ Reject", fmt.Sprintf("reject:%d", confID)),
			},
			{
				tgbotapi.NewInlineKeyboardButtonData("✏️ Edit", fmt.Sprintf("edit:%d", confID)),
				tgbotapi.NewInlineKeyboardButtonData("🗑 Delete", fmt.Sprintf("delete:%d", confID)),
			},
			{
				tgbotapi.NewInlineKeyboardButtonData("⏳ Schedule 1h", fmt.Sprintf("schedule:%d", confID)),
				tgbotapi.NewInlineKeyboardButtonData("⛔ Ban Sender", fmt.Sprintf("ban:%d", uid)),
			},
			{
				tgbotapi.NewInlineKeyboardButtonData("⬅️ Prev", fmt.Sprintf("page:%d", confID-1)),
				tgbotapi.NewInlineKeyboardButtonData("➡️ Next", fmt.Sprintf("page:%d", confID+1)),
			},
		},
	}
}

// ---------------- ADMIN CALLBACKS ----------------

func handleCallback(cb *tgbotapi.CallbackQuery) {
	if cb.From.ID != AdminID {
		bot.AnswerCallbackQuery(tgbotapi.NewCallback(cb.ID, "⛔ Admin only"))
		return
	}

	data := cb.Data

	switch {
	case strings.HasPrefix(data, "approve:"):
		var confID, uid int64
		fmt.Sscanf(data, "approve:%d:%d", &confID, &uid)
		publishConfession(confID)
		editAdminMsg(cb.Message, "✅ Approved & Posted")

	case strings.HasPrefix(data, "reject:"):
		confID := parseID(data)
		db.Exec("UPDATE confessions SET status='rejected' WHERE id=?", confID)
		editAdminMsg(cb.Message, "❌ Rejected")

	case strings.HasPrefix(data, "edit:"):
		confID := parseID(data)
		bot.Send(tgbotapi.NewMessage(cb.Message.Chat.ID, fmt.Sprintf("✏️ Edit confession ID %d via /edit command.", confID)))

	case strings.HasPrefix(data, "delete:"):
		confID := parseID(data)
		deleteFromChannel(confID)
		editAdminMsg(cb.Message, "🗑 Deleted")

	case strings.HasPrefix(data, "schedule:"):
		confID := parseID(data)
		scheduleConfession(confID, nowUnix()+3600)
		editAdminMsg(cb.Message, "⏳ Scheduled 1 hour later")

	case strings.HasPrefix(data, "ban:"):
		uid := parseID(data)
		db.Exec("UPDATE users SET banned=1 WHERE user_id=?", uid)
		editAdminMsg(cb.Message, "⛔ User banned")

	case strings.HasPrefix(data, "page:"):
		id := parseID(data)
		showConfessionPage(cb.Message.Chat.ID, id)
	}

	bot.AnswerCallbackQuery(tgbotapi.NewCallback(cb.ID, "✔ Done"))
}

// ---------------- HELPERS ----------------

func parseID(data string) int64 {
	var id int64
	fmt.Sscanf(data[strings.Index(data, ":")+1:], "%d", &id)
	return id
}

func editAdminMsg(msg *tgbotapi.Message, status string) {
	bot.EditMessageText(tgbotapi.NewEditMessageText(msg.Chat.ID, msg.MessageID, status))
}
