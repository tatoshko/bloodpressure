package handlerLog

import (
    "fmt"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
    "regexp"
    "strconv"
    "useful.team/bloodpressure/m/bot/core"
)

const (
    Short = `^(?P<up>\d{2,3})/(?P<down>\d{2,3})/(?P<pulse>\d{2,3})$`
)

func LogShort(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
    logger := getLogger("LogShort")

    userID := update.Message.From.ID
    chatID := update.Message.Chat.ID

    params := core.GetParams(Short, update.Message.Text)

    if len(params) == 3 {
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
                msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Добавленно: %d/%d пульс %d", up, down, pulse))
                if _, err := bot.Send(msg); err != nil {
                    logger(err.Error())
                }
            }
        }
    }
}

func Check(message string) bool {
    return regexp.MustCompile(Short).MatchString(message)
}
