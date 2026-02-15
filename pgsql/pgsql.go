package pgsql

import (
    "database/sql"
    "fmt"
    _ "github.com/lib/pq"
)

var (
    db           *sql.DB
    storedConfig Config
)

func Init(config Config) *sql.DB {
    var err error
    logger := getLogger("INIT")

    if db != nil {
        return db
    }

    storedConfig = config
    psqlInfo := fmt.Sprintf(
        "host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
        config.Host, config.Port, config.User, config.Password, config.DBName,
    )

    logger("Trying to connect to: ", psqlInfo)

    db, err = sql.Open("postgres", psqlInfo)
    if err != nil {
        panic(err)
    }

    err = db.Ping()
    if err != nil {
        panic(err)
    }

    return db
}

func GetClient() *sql.DB {
    if db == nil {
        Init(storedConfig)
    }

    return db
}
