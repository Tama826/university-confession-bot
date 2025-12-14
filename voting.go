package main

import (
        "fmt"

        tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// ---------------- VOTING ----------------

func voteKeyboard(confID int64, up int, down int) *tgbotapi.InlineKeyboardMarkup {
        return &tgbotapi.InlineKeyboardMarkup{
                InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{
                        {
                                tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("👍 %d", up), fmt.Sprintf("vote:%d:1", confID)),
                                tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("👎 %d", down), fmt.Sprintf("vote:%d:-1", confID)),
                        },
                },
        }
}

func handleVote(cb *tgbotapi.CallbackQuery) {
        var confID int64
        var vote int
        fmt.Sscanf(cb.Data, "vote:%d:%d", &confID, &vote)
        uid := cb.From.ID

        var prev int
        err := db.QueryRow("SELECT vote FROM votes WHERE confession_id=? AND user_id=?", confID, uid).Scan(&prev)

        if err == nil {
                if prev == vote {
                        callback := tgbotapi.NewCallback(cb.ID, "❌ You already voted")
                        bot.Request(callback)
                        return
                }
                db.Exec("UPDATE votes SET vote=? WHERE confession_id=? AND user_id=?", vote, confID, uid)
                updateVoteCounts(confID, -prev, vote)
        } else {
                db.Exec("INSERT INTO votes(confession_id,user_id,vote) VALUES(?,?,?)", confID, uid, vote)
                updateVoteCounts(confID, 0, vote)
        }

        refreshVoteMessage(cb.Message, confID)
        callback := tgbotapi.NewCallback(cb.ID, "✔ Vote recorded")
        bot.Request(callback)
}

func updateVoteCounts(confID int64, remove int, add int) {
        if remove == 1 {
                db.Exec("UPDATE confessions SET upvotes=upvotes-1 WHERE id=?", confID)
        }
        if remove == -1 {
                db.Exec("UPDATE confessions SET downvotes=downvotes-1 WHERE id=?", confID)
        }
        if add == 1 {
                db.Exec("UPDATE confessions SET upvotes=upvotes+1 WHERE id=?", confID)
        }
        if add == -1 {
                db.Exec("UPDATE confessions SET downvotes=downvotes+1 WHERE id=?", confID)
        }
}

func refreshVoteMessage(msg *tgbotapi.Message, confID int64) {
        var up, down int
        db.QueryRow("SELECT upvotes, downvotes FROM confessions WHERE id=?", confID).Scan(&up, &down)
        edit := tgbotapi.NewEditMessageReplyMarkup(msg.Chat.ID, msg.MessageID, *voteKeyboard(confID, up, down))
        bot.Send(edit)
}

func resetVotes(confID int64) {
        db.Exec("DELETE FROM votes WHERE confession_id=?", confID)
        db.Exec("UPDATE confessions SET upvotes=0, downvotes=0 WHERE id=?", confID)
}
