package handlerLog

import "useful.team/bloodpressure/m/pgsql"

type LogService struct {
    UserId int64
}

func NewLogService(userId int64) *LogService {
    return &LogService{UserId: userId}
}

func (ls LogService) Add(up, down, pulse int) (err error) {
    pg := pgsql.GetClient()
    q := `insert into log (user_uuid, up, down, pulse) VALUES ($1, $2, $3, $4)`

    _, err = pg.Exec(q, ls.UserId, up, down, pulse)
    return err
}
