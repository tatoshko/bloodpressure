package bot

import "fmt"

type Config struct {
    Name        string `env:"NAME"`
    Token       string `env:"TOKEN", required`
    Hook        string `env:"HOOK", required`
    APIEndpoint string `env:"API_ENDPOINT"`
}

func (c Config) ToString() string {
    var token string
    if len(c.Token) < 5 {
        token = "empty"
    } else {
        token = fmt.Sprintf("%s..%s", c.Token[0:2], c.Token[len(c.Token)-3:])
    }
    return fmt.Sprintf("BOT:\n  Name: %s\n  Token: %s\n  Hook: %s\n  API: %s\n", c.Name, token, c.Hook, c.APIEndpoint)
}
