package handlerLog

import (
    "database/sql"
    "sort"
    "time"
    "useful.team/bloodpressure/m/bot/core"
    "useful.team/bloodpressure/m/pgsql"
)

type LogRecord struct {
    UUID      string    `json:"UUID,omitempty"`
    UserUUID  string    `json:"user_uuid,omitempty"`
    Up        int       `json:"up,omitempty"`
    Down      int       `json:"down,omitempty"`
    Pulse     int       `json:"pulse,omitempty"`
    CreatedAt time.Time `json:"created_at"`
}

func (lr *LogRecord) Score() int {
    return lr.Up + lr.Down
}

type LogStat struct {
    HigherPressure *LogRecord `json:"higher_pressure,omitempty"`
    LowerPressure  *LogRecord `json:"lower_pressure,omitempty"`
    HigherPulse    *LogRecord `json:"higher_pulse,omitempty"`
    LowerPulse     *LogRecord `json:"lower_pulse,omitempty"`
}

type LogService struct {
    User *core.User
}

func NewLogService(user *core.User) *LogService {
    return &LogService{User: user}
}

func (ls *LogService) Add(up, down, pulse int) (err error) {
    pg := pgsql.GetClient()
    q := `insert into log (user_uuid, up, down, pulse) VALUES ($1, $2, $3, $4)`

    _, err = pg.Exec(q, ls.User.UUID, up, down, pulse)
    return err
}

func (ls *LogService) FindLastMonthToNow() (logRecords []*LogRecord, err error) {
    pg := pgsql.GetClient()
    q := `select uuid, user_uuid, up, down, pulse, created_at from log 
        where user_uuid = $1 and created_at >= date_trunc('month', current_date - interval '1 month' )
        order by created_at asc`

    var rows *sql.Rows
    if rows, err = pg.Query(q, ls.User.UUID); err != nil {
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        lr := &LogRecord{}
        if err = rows.Scan(&lr.UUID, &lr.UserUUID, &lr.Up, &lr.Down, &lr.Pulse, &lr.CreatedAt); err != nil {
            return nil, err
        } else {
            logRecords = append(logRecords, lr)
        }
    }

    return logRecords, nil
}

func (ls *LogService) FindLastYear() (logRecords []*LogRecord, err error) {
    pg := pgsql.GetClient()
    q := `select uuid, user_uuid, up, down, pulse, created_at from log
        where user_uuid = $1 and created_at >= date_trunc('year', current_date)
        order by created_at asc`

    var rows *sql.Rows
    if rows, err = pg.Query(q, ls.User.UUID); err != nil {
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        lr := &LogRecord{}
        if err = rows.Scan(&lr.UUID, &lr.UserUUID, &lr.Up, &lr.Down, &lr.Pulse, &lr.CreatedAt); err != nil {
            return nil, err
        } else {
            logRecords = append(logRecords, lr)
        }
    }

    return logRecords, nil
}

func (ls *LogService) FindStatistic() (stat *LogStat, err error) {
    pg := pgsql.GetClient()
    q := `select name, uuid, user_uuid, up, down, pulse, created_at from (
        (
            select 'lower pressure' name, *, (up + down) sum from log
            where user_uuid = $1 
            order by sum limit 1
        ) union (
            select 'higher pressure' name, *, (up + down) sum from log
            where user_uuid = $1
            order by sum desc limit 1
        ) union (
            select 'lower pulse' name, *, (up + down) sum from log
            where user_uuid = $1
            order by pulse limit 1
        ) union (
            select 'higher pulse' name, *, (up + down) sum from log
            where user_uuid = $1
            order by pulse desc limit 1
        )
    );`

    var rows *sql.Rows

    if rows, err = pg.Query(q, ls.User.UUID); err != nil {
        return nil, err
    }
    defer rows.Close()

    stat = &LogStat{}

    for rows.Next() {
        var name string
        var record *LogRecord

        if err = rows.Scan(
            &name,
            &record.UUID,
            &record.UserUUID,
            &record.Up,
            &record.Down,
            &record.Pulse,
            &record.CreatedAt,
        ); err != nil {
            return nil, err
        }

        switch name {
        case "lower_pressure":
            stat.LowerPressure = record
        case "higher_pressure":
            stat.HigherPressure = record
        case "lower_pulse":
            stat.LowerPulse = record
        case "higher_pulse":
            stat.HigherPulse = record
        }
    }

    return stat, nil
}

func (ls *LogService) ComputeMedian(records []*LogRecord) *LogRecord {
    l := len(records)

    if l <= 0 {
        return nil
    }

    sort.Slice(records, func(i, j int) bool {
        return records[i].Score() > records[j].Score()
    })

    if l%2 != 0 {
        return records[((l+1)/2)-1]
    }

    a := records[(l/2)-1]
    b := records[l/2]

    if a.Score() > b.Score() {
        return a
    } else {
        if a.Up > b.Up {
            return a
        }
        return b
    }
}
