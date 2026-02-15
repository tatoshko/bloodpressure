package handlerLog

import (
    "fmt"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
    "useful.team/bloodpressure/m/bot/core"
)

func Stat(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
    var err error
    logger := getLogger("GetStat")

    userID := update.Message.From.ID
    chatID := update.Message.Chat.ID

    userService := core.NewUserService()
    var user *core.User

    if user, err = userService.FindById(userID); err != nil {
        msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Пользователь не найден: %s", err.Error()))

        if _, err := bot.Send(msg); err != nil {
            logger(err.Error())
        }
        return
    }

    logService := NewLogService(user)

    var logRecords []*LogRecord
    if logRecords, err = logService.GetStat(30); err != nil {
        msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Не смог достать записи: %s", err.Error()))

        if _, err := bot.Send(msg); err != nil {
            logger(err.Error())
        }
        return
    }

    var statMessage string
    for _, record := range logRecords {
        statMessage += fmt.Sprintf("[%s] Давление: <b>%d/%d</b>, Пульс: <b>%d</b>\n", record.CreatedAt.Format("02 Jan 15:04"), record.Up, record.Down, record.Pulse)
    }

    msg := tgbotapi.NewMessage(chatID, statMessage)
    msg.ParseMode = tgbotapi.ModeHTML

    if _, err := bot.Send(msg); err != nil {
        logger(err.Error())
    }
}
