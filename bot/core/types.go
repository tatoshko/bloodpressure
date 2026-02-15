package core

import tba "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type Handler func(api *tba.BotAPI, update tba.Update)
