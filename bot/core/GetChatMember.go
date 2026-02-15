package core

import (
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
    "log"
)

func GetChatMember(bot *tgbotapi.BotAPI, chatID int64, userID int64) (member tgbotapi.ChatMember, err error) {
    config := tgbotapi.ChatConfigWithUser{
        ChatID: chatID,
        UserID: userID,
    }
    if member, err = bot.GetChatMember(tgbotapi.GetChatMemberConfig{ChatConfigWithUser: config}); err != nil {
        log.Printf("HandleReg GetChatMember err %s", err.Error())
    }

    return
}
