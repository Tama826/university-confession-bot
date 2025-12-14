package main

import (
        "encoding/json"
        "log"
        "net/http"
        "strings"

        tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func webhookHandler(w http.ResponseWriter, r *http.Request) {
        log.Printf("Received request: %s %s", r.Method, r.URL.Path)

        if r.Method != http.MethodPost {
                w.WriteHeader(http.StatusOK)
                return
        }

        var update tgbotapi.Update
        if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
                log.Printf("Error decoding update: %v", err)
                w.WriteHeader(http.StatusBadRequest)
                return
        }

        log.Printf("Update received: %+v", update.UpdateID)

        if update.Message != nil {
                log.Printf("Message from %d: %s", update.Message.From.ID, update.Message.Text)
                handleMessage(update.Message)
        }

        if update.CallbackQuery != nil {
                log.Printf("Callback from %d: %s", update.CallbackQuery.From.ID, update.CallbackQuery.Data)
                if strings.HasPrefix(update.CallbackQuery.Data, "vote:") {
                        handleVote(update.CallbackQuery)
                } else {
                        handleCallback(update.CallbackQuery)
                }
        }

        w.WriteHeader(http.StatusOK)
}
