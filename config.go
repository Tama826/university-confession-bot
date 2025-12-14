package main

import "time"

// ---------- BOT CONFIGURATION ----------
var (
	AdminID      int64 = 123456789       // Telegram Owner ID
	AdminGroupID int64 = -1001111111111  // Group ID for admin moderation
	ChannelID    int64 = -1002222222222  // Channel to post confessions
)

// ---------- BLIND CHAT ----------
var (
	WaitingUser  int64
	WaitingSince int64
	Pairs        = make(map[int64]int64)
	Reports      = make(map[int64]int)
)

// ---------- SCHEDULER ----------
var ScheduleCheckInterval = 20 * time.Second

// ---------- MODERATION ----------
var BannedWords = []string{
	"rape", "kill", "terror", "bomb",
}

// ---------- VOTING ----------
var MaxVotePerUser = 1
