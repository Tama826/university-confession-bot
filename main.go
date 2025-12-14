package main

import (
        "log"
        "net/http"
        "os"

        tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
        bot *tgbotapi.BotAPI
)

func main() {
        var err error

        // Init bot
        bot, err = tgbotapi.NewBotAPI(os.Getenv("BOT_TOKEN"))
        if err != nil {
                log.Fatal(err)
        }

        log.Printf("Authorized as @%s", bot.Self.UserName)

        // Init database
        initDatabase()

        // Start background workers
        startScheduler()

        // Webhook
        http.HandleFunc("/", webhookHandler)

        log.Println("Bot running (webhook mode)")
        log.Fatal(http.ListenAndServe("0.0.0.0:5000", nil))
}
