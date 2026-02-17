package handlerStart

import (
    "fmt"
    tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
    "useful.team/bloodpressure/m/bot/core"
)

func Start(bot *tgbotapi.BotAPI, update tgbotapi.Update) {
    var err error

    logger := getLogger("Add new user")
    userId := update.Message.From.ID
    chatId := update.Message.Chat.ID

    userService := core.NewUserService()

    if exist, err := userService.CheckExist(userId); err != nil {
        logger(err.Error())
        return
    } else if exist {
        msg := tgbotapi.NewMessage(chatId, "Вы уже зарегистрированы")
        if _, err = bot.Send(msg); err != nil {
            logger(err.Error())
        }
        return
    }

    if err = userService.Add(userId); err != nil {
        logger(err.Error())

        msg := tgbotapi.NewMessage(chatId, fmt.Sprintf("Ошибка добавления нового пользователя: %s", err.Error()))
        if _, err = bot.Send(msg); err != nil {
            logger(err.Error())
        }

        return
    } else {
        msg := tgbotapi.NewMessage(chatId, "Добро пожаловать! Что бы добавить в журнал запись, просто отправьте сообщение вида:\n<code>120 80 60</code>\nГде 120 верхнее давление, 80 нижнее, 60 пульс")
        msg.ParseMode = tgbotapi.ModeHTML

        if _, err = bot.Send(msg); err != nil {
            logger(err.Error())
        }
    }
}
