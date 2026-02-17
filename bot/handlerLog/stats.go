package handlerLog

import (
    "bytes"
    "fmt"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
    "github.com/xuri/excelize/v2"
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
    medianPressure := logService.FindMedian(logRecords)

    statMessage += "За последнее время\n"
    statMessage += fmt.Sprintf(
        "Самое высокое давление:\n<b>%d/%d</b> при пульсе: %d, было %s\n\n",
        highestPressure.Up,
        highestPressure.Down,
        highestPressure.Pulse,
        highestPressure.CreatedAt.Format("02 Jan 15:04"),
    )
    statMessage += fmt.Sprintf(
        "Самый высокий пульс:\n%s\n<b>%d</b> при давлении: %d/%d",
        highestPulse.CreatedAt.Format("02 Jan 15:04"),
        highestPulse.Pulse,
        highestPressure.Up,
        highestPressure.Down,
    )
    statMessage += fmt.Sprintf(
        "Медианное значение: %d/%d",
        medianPressure.Up,
        medianPressure.Down,
    )

    msg := tgbotapi.NewMessage(chatID, statMessage)
    msg.ParseMode = tgbotapi.ModeHTML

    if _, err := bot.Send(msg); err != nil {
        logger(err.Error())
    }

    // Send excel with log data
    f := excelize.NewFile()
    defer func() {
        if err := f.Close(); err != nil {
            logger(err.Error())
        }
    }()

    f.SetSheetName("Sheet1", "Log")

    f.SetCellStr("Log", "A1", "Date")
    f.SetCellStr("Log", "B1", "Up")
    f.SetCellStr("Log", "C1", "Down")
    f.SetCellStr("Log", "D1", "Pulse")

    for i, record := range logRecords {
        f.SetCellValue("Log", fmt.Sprintf("A%d", i+2), record.CreatedAt)
        f.SetCellInt("Log", fmt.Sprintf("B%d", i+2), int64(record.Up))
        f.SetCellInt("Log", fmt.Sprintf("C%d", i+2), int64(record.Down))
        f.SetCellInt("Log", fmt.Sprintf("D%d", i+2), int64(record.Pulse))
    }

    var buf *bytes.Buffer
    if buf, err = f.WriteToBuffer(); err != nil {
        logger(err.Error())
        return
    }

    file := tgbotapi.FileBytes{Name: "log.xlsx", Bytes: buf.Bytes()}
    msgXlsx := tgbotapi.NewDocument(chatID, file)

    if _, err := bot.Send(msgXlsx); err != nil {
        logger(err.Error())
    }
}
