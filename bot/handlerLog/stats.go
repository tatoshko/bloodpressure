package handlerLog

import (
    "fmt"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
    "useful.team/bloodpressure/m/bot/core"
)

func Stat(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
    var err error
    logger := getLogger("GetLast")

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
    lastRecordsCount := 30

    var logRecords []*LogRecord
    if logRecords, err = logService.GetLast(lastRecordsCount); err != nil {
        msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Не смог достать записи: %s", err.Error()))

        if _, err := bot.Send(msg); err != nil {
            logger(err.Error())
        }
        return
    }

    // Send Summary
    var statMessage string

    highestPressure := logService.FindHighestPressure(logRecords)
    highestPulse := logService.FindHighestPulse(logRecords)
    medianPressure := logService.FindMedian(logRecords[:])

    statMessage += "За последнее время\n\n"
    statMessage += fmt.Sprintf(
        "Самое высокое давление:\n<b>%d/%d</b> при пульсе <b>%d</b> было %s\n\n",
        highestPressure.Up,
        highestPressure.Down,
        highestPressure.Pulse,
        highestPressure.CreatedAt.Format("02 Jan 15:04"),
    )
    statMessage += fmt.Sprintf(
        "Самый высокий пульс:\n<b>%d</b> при давлении <b>%d/%d</b> был %s\n\n",
        highestPulse.Pulse,
        highestPulse.Up,
        highestPulse.Down,
        highestPulse.CreatedAt.Format("02 Jan 15:04"),
    )
    statMessage += fmt.Sprintf(
        "Медианное значение:\n<b>%d/%d</b>",
        medianPressure.Up,
        medianPressure.Down,
    )

    msg := tgbotapi.NewMessage(chatID, statMessage)
    msg.ParseMode = tgbotapi.ModeHTML

    if _, err := bot.Send(msg); err != nil {
        logger(err.Error())
    }
}
