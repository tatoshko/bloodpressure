package main

import (
    "fmt"
    "github.com/caarlos0/env/v11"
    "log"
    "useful.team/bloodpressure/m/bot"
    "useful.team/bloodpressure/m/httpServer"
    "useful.team/bloodpressure/m/pgsql"
)

var err error

type Config struct {
    Bot    bot.Config        `envPrefix:"BOT_"`
    Pg     pgsql.Config      `envPrefix:"PG_"`
    Server httpServer.Config `envPrefix:"SERVER_"`
}

func (c Config) ToString() string {
    return fmt.Sprintf("%s%s%s\n", c.Server.ToString(), c.Pg.ToString(), c.Bot.ToString())
}

func main() {
    var config Config
    if err := env.Parse(&config); err != nil {
        log.Fatalln(err.Error())
    }

    log.Printf("Config loaded: \n%s\n", config.ToString())

    go httpServer.Init(config.Server)
    go pgsql.Init(config.Pg)
    bot.Start(config.Bot)
}
