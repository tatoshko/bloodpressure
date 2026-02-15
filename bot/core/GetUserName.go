package core

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func GetUserName(user *tgbotapi.User, mention bool) string {
    if user.UserName != "" {
        if mention {
            return "@" + user.UserName
        }
        return user.UserName
    } else if user.FirstName != "" {
        return user.FirstName + " " + user.LastName
    } else if user.LastName != "" {
        return user.LastName
    }

    return "Хер_знает_кто_такой"
}
