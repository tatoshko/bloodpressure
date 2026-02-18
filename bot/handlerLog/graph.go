package handlerLog

import (
    "fmt"
    "github.com/go-echarts/go-echarts/v2/charts"
    "github.com/go-echarts/go-echarts/v2/opts"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
    "time"
    "useful.team/bloodpressure/m/bot/core"
)

func Graph(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
    var err error
    logger := getLogger("XLSX")

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
    if logRecords, err = logService.FindLastMonthToNow(); err != nil {
        msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Не смог достать записи: %s", err.Error()))

        if _, err := bot.Send(msg); err != nil {
            logger(err.Error())
        }
        return
    }

    var dates []time.Time
    var ups, downs, pulses []opts.LineData
    for _, record := range logRecords {
        dates = append(dates, record.CreatedAt)
        ups = append(ups, opts.LineData{Value: record.Up})
        downs = append(downs, opts.LineData{Value: record.Down})
        pulses = append(pulses, opts.LineData{Value: record.Pulse})
    }

    line := charts.NewLine()
    line.SetXAxis(dates).
        AddSeries("Ups", ups).
        AddSeries("Downs", downs).
        AddSeries("Pulses", pulses)

    file := tgbotapi.FileBytes{Bytes: line.RenderContent(), Name: "Chart"}
    msg := tgbotapi.NewDocument(chatID, file)

    if _, err := bot.Send(msg); err != nil {
        logger(err.Error())
    }

}
