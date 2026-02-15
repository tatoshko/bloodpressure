package handlerNew

import (
    "useful.team/bloodpressure/m/pgsql"
)

type UserService struct {
}

func NewUserService() *UserService {
    return &UserService{}
}

func (us *UserService) Add(telegramId int64) (err error) {
    pg := pgsql.GetClient()
    q := `insert into user (telegram_id) values ($1)`

    _, err = pg.Exec(q, telegramId)
    return
}
