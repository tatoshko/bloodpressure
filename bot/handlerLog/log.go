package handlerLog

import (
    "fmt"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
    "strconv"
    "useful.team/bloodpressure/m/bot/core"
)

func Log(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
    logger := getLogger("Log")

    userID := update.Message.From.ID
    chatID := update.Message.Chat.ID

    params := core.GetParams(`^(?P<up>\d+)\D+(?P<down>\d+)\D+(?P<pulse>\d+)$`, update.Message.Text)

    if len(params) == 3 {
        //logger(fmt.Sprintf("UserID: %d, Params: %v", userID, params))

        userService := core.NewUserService()

        if user, err := userService.FindById(userID); err != nil {
            msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Пользователь не найден: %s", err.Error()))

            if _, err := bot.Send(msg); err != nil {
                logger(err.Error())
            }
            return
        } else {
            up, _ := strconv.Atoi(params["up"])
            down, _ := strconv.Atoi(params["down"])
            pulse, _ := strconv.Atoi(params["pulse"])

            logService := NewLogService(user)
            if err := logService.Add(up, down, pulse); err != nil {
                logger(err.Error())
            } else {
                msg := tgbotapi.NewMessage(chatID, "OK")
                if _, err := bot.Send(msg); err != nil {
                    logger(err.Error())
                }
            }
        }
    }
}
