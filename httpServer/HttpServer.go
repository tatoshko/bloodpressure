package httpServer

import (
    "fmt"
    "log"
    "net/http"
)

type Config struct {
    Host string `env:"HOST"`
    Port int    `env:"PORT"`
}

func (c Config) ToString() string {
    return fmt.Sprintf("HTTP SERVER:\n  Host: %s\n  Port: %d\n", c.Host, c.Port)
}

func Init(config Config) {
    logger := getLogger("INIT")
    addr := fmt.Sprintf("%s:%d", config.Host, config.Port)

    logger(fmt.Sprintf("Trying bind to: [%s]", addr))

    if err := http.ListenAndServe(addr, nil); err != nil {
        log.Fatalln(err.Error())
    }

    logger("HTTP server is ready")
}
