package handlerLog

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func Log(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
    userID := update.Message.From.ID

}
