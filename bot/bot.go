package bot

import (
    "context"
    tba "github.com/go-telegram-bot-api/telegram-bot-api/v5"
    "log"
    "net"
    "net/http"
    "strings"
    "time"
    "useful.team/bloodpressure/m/bot/callbacks"
    "useful.team/bloodpressure/m/bot/core"
    "useful.team/bloodpressure/m/bot/handlerLog"
    "useful.team/bloodpressure/m/bot/handlerStart"
)

var (
    Commands = make(map[string]core.Handler)
)

func Start(config Config) {
    transport := &http.Transport{
        DialContext: (&net.Dialer{
            Timeout:   30 * time.Second,
            KeepAlive: 30 * time.Second,
        }).DialContext,
        TLSHandshakeTimeout: 10 * time.Second,
        // Force using specific IP for api.telegram.org
    }

    // Подменяем разрешение имени
    transport.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
        log.Printf("ADDR: %s", addr)
        if addr == "api.telegram.org" {
            // Используем один из известных IP Telegram
            // Можно попробовать несколько: 149.154.167.220, 149.154.167.221, 149.154.167.51
            addr = "149.154.167.220"
        }
        // Для всех остальных адресов используем стандартный Dialer
        return (&net.Dialer{
            Timeout:   30 * time.Second,
            KeepAlive: 30 * time.Second,
        }).DialContext(ctx, network, addr)
    }

    hc := &http.Client{
        Transport: transport,
        Timeout:   30 * time.Second,
    }

    if API, err := tba.NewBotAPIWithClient(config.Token, config.APIEndpoint, hc); err == nil {

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
                    if handlerLog.Check(update.Message.Text) {
                        go handlerLog.LogShort(API, update)
                    }
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
    Commands["xlsx"] = handlerLog.Xlsx
    Commands["graph"] = handlerLog.Graph
}
