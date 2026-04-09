package pgsql

import "fmt"

type Config struct {
    Host     string `env:"HOST"`
    Port     int    `env:"PORT"`
    User     string `env:"USER"`
    Password string `env:"PASSWORD"`
    Name     string `env:"NAME"`
}

func (c Config) ToString() string {
    var password string
    if len(c.Password) < 3 {
        password = "empty"
    } else {
        password = fmt.Sprintf("..%s", c.Password[len(c.Password)-3:])
    }

    return fmt.Sprintf(
        "PG:\n  Host: %s\n  Port: %d\n  User: %s\n  Password: %s\n  Name: %s\n",
        c.Host, c.Port, c.User, password, c.Name,
    )
}
