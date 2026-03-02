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

    // Send Summary
    var statMessage string

    var stat *LogStat
    if stat, err = logService.FindStatistic(); err != nil {
        logger(err.Error())
        return
    }

    statMessage += "За последний месяц\n"

    statMessage += fmt.Sprintf(
        "%d/%d - %d/%d\n\n",
        stat.LowerPressure.Up, stat.LowerPressure.Down,
        stat.HigherPressure.Up, stat.HigherPressure.Down,
    )

    var logRecords []*LogRecord
    if logRecords, err = logService.FindLastMonthToNow(); err != nil {
        msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Не смог достать записи: %s", err.Error()))

        if _, err := bot.Send(msg); err != nil {
            logger(err.Error())
        }
        return
    }
    medianPressure := logService.ComputePressureMedian(logRecords)

    statMessage += fmt.Sprintf(
        "Медианные:\n<b>%d/%d</b>",
        medianPressure.Up,
        medianPressure.Down,
    )

    msg := tgbotapi.NewMessage(chatID, statMessage)
    msg.ParseMode = tgbotapi.ModeHTML

    if _, err := bot.Send(msg); err != nil {
        logger(err.Error())
    }
}
