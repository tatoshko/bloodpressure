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

    logger(fmt.Sprintf("%v", stat.HigherPressure))

    statMessage += fmt.Sprintf(
        "Самое высокое давление:\n<b>%d/%d</b> при пульсе <b>%d</b> было %s\n\n",
        stat.HigherPressure.Up,
        stat.HigherPressure.Down,
        stat.HigherPressure.Pulse,
        stat.HigherPressure.CreatedAt.Format("02 Jan 15:04"),
    )
    statMessage += fmt.Sprintf(
        "Самое низкое давление:\n<b>%d/%d</b> при пульсе <b>%d</b> было %s\n\n",
        stat.LowerPressure.Up,
        stat.LowerPressure.Down,
        stat.LowerPressure.Pulse,
        stat.LowerPressure.CreatedAt.Format("02 Jan 15:04"),
    )

    statMessage += fmt.Sprintf(
        "Самый высокий пульс:\n<b>%d</b> при давлении <b>%d/%d</b> был %s\n\n",
        stat.HigherPulse.Pulse,
        stat.HigherPulse.Up,
        stat.HigherPulse.Down,
        stat.HigherPulse.CreatedAt.Format("02 Jan 15:04"),
    )
    statMessage += fmt.Sprintf(
        "Самый низкий пульс:\n<b>%d</b> при давлении <b>%d/%d</b> был %s\n\n",
        stat.LowerPulse.Pulse,
        stat.LowerPulse.Up,
        stat.LowerPulse.Down,
        stat.LowerPulse.CreatedAt.Format("02 Jan 15:04"),
    )

    var logRecords []*LogRecord
    if logRecords, err = logService.FindLastMonthToNow(); err != nil {
        msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Не смог достать записи: %s", err.Error()))

        if _, err := bot.Send(msg); err != nil {
            logger(err.Error())
        }
        return
    }
    medianPressure := logService.ComputeMedian(logRecords[:])

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
