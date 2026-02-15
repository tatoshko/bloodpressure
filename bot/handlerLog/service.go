package handlerLog

import (
    "useful.team/bloodpressure/m/bot/core"
    "useful.team/bloodpressure/m/pgsql"
)

type LogService struct {
    User *core.User
}

func NewLogService(user *core.User) *LogService {
    return &LogService{User: user}
}

func (ls LogService) Add(up, down, pulse int) (err error) {
    pg := pgsql.GetClient()
    q := `insert into log (user_uuid, up, down, pulse) VALUES ($1, $2, $3, $4)`

    _, err = pg.Exec(q, ls.User.UUID, up, down, pulse)
    return err
}
