package handlerStart

import (
    "fmt"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Start(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
    var err error

    logger := getLogger("Add new user")
    userId := update.Message.From.ID
    chatId := update.Message.Chat.ID

    userService := NewUserService()

    if err = userService.Add(userId); err != nil {
        logger(err.Error())

        msg := tgbotapi.NewMessage(chatId, fmt.Sprintf("Ошибка добавления нового пользователя: %s", err.Error()))
        if _, err = bot.Send(msg); err != nil {
            logger(err.Error())
        }

        return
    } else {
        msg := tgbotapi.NewMessage(chatId, "Добро пожаловать! Что бы добавить в журнал запись, просто отправьте сообщение вида 120 80 60 - где 120 верхнее давление, 80 нижнее, 60 пульс")

        if _, err = bot.Send(msg); err != nil {
            logger(err.Error())
        }
    }
}
