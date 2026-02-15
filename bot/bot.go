package bot

import (
    tba "github.com/go-telegram-bot-api/telegram-bot-api/v5"
    "log"
    "strings"
    "useful.team/bloodpressure/m/bot/callbacks"
    "useful.team/bloodpressure/m/bot/core"
    "useful.team/bloodpressure/m/bot/handlerLog"
    "useful.team/bloodpressure/m/bot/handlerStart"
)

var (
    Commands = make(map[string]core.Handler)
)

func Start(config Config) {
    if API, err := tba.NewBotAPI(config.Token); err == nil {

        wh, _ := tba.NewWebhook(config.Hook + "/" + config.Token)
        if _, err := API.Request(wh); err != nil {
            log.Printf("SetHoook error %s\n", err.Error())
        }

        API.Debug = false

        registerCommands()

        updates := API.ListenForWebhook("/" + API.Token)

        for update := range updates {
            if update.Message != nil {
                message := update.Message

                direct := int64(message.From.ID) == message.Chat.ID
                tagMe := strings.Index(message.CommandWithAt(), config.Name) != -1

                if message.IsCommand() && (tagMe || direct) {
                    if handler, found := Commands[message.Command()]; found {
                        log.Printf(
                            "MessageID: '%d', Command: '%s', Data: '%s', From: '%d'\n",
                            message.MessageID, message.Command(), message.CommandArguments(), message.From.ID,
                        )
                        go handler(API, update)
                    }
                } else {
                    go handlerLog.Log(API, update)
                }
            } else if update.CallbackQuery != nil {
                data := update.CallbackQuery.Data

                log.Printf("CallBackQuery %s", data)

                var handlerID string
                if strings.HasPrefix(data, "/") {
                    parts := strings.SplitN(data, " ", 2)
                    handlerID = strings.TrimPrefix(parts[0], "/")
                } else {
                    handlerID = data
                }

                if handler, found := callbacks.GetHandler(handlerID); found {
                    handler(API, update)
                }
            }
        }
    } else {
        log.Fatalf("NewAPIBot error %s\n", err.Error())
    }
}

func registerCommands() {
    Commands["start"] = handlerStart.Start
    Commands["stat"] = handlerLog.Stat
}
