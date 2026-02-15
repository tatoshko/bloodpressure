package handlerLog

import (
    "fmt"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
    "useful.team/bloodpressure/m/bot/core"
)

func Log(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
    logger := getLogger("Log")

    userID := update.Message.From.ID
    params := core.GetParams(`^(?P<up>\d+)\D+(?P<down>\d+)\D+(?P<pulse>\d+)$`, update.Message.Text)

    //logService := NewLogService(userID)

    logger(fmt.Sprintf("UserID: %d, Params: %v", userID, params))
}
