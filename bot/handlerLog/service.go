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

func (ls *LogService) GetLast(limit int) (logRecords []*LogRecord, err error) {
    pg := pgsql.GetClient()
    q := `select uuid, user_uuid, up, down, pulse, created_at from log where user_uuid = $1 order by created_at asc limit $2`

    var rows *sql.Rows
    if rows, err = pg.Query(q, ls.User.UUID, limit); err != nil {
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

func (ls *LogService) FindHighestPressure(records []*LogRecord) *LogRecord {
    if len(records) <= 0 {
        return nil
    }

    max := records[0]
    maxSum := max.Up + max.Down

    for _, record := range records[1:] {
        recordSum := record.Up + record.Down
        if recordSum > maxSum {
            max = record
            maxSum = recordSum
        }
    }

    return max
}

func (ls *LogService) FindHighestPulse(records []*LogRecord) *LogRecord {
    if len(records) <= 0 {
        return nil
    }

    max := records[0]

    for _, record := range records[1:] {
        if record.Pulse > max.Pulse {
            max = record
        }
    }

    return max
}

func (ls *LogService) FindMedian(records []*LogRecord) *LogRecord {
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

    if (a.Up + a.Down) > (b.Up + b.Down) {
        return a
    } else {
        return b
    }
}
