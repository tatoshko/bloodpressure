package core

import (
    "database/sql"
    "time"
    "useful.team/bloodpressure/m/pgsql"
)

type User struct {
    UUID       string    `json:"UUID,omitempty"`
    TelegramId int64     `json:"telegram_id,omitempty"`
    CreatedAt  time.Time `json:"created_at"`
}

type UserService struct {
}

func NewUserService() *UserService {
    return &UserService{}
}

func (us UserService) FindById(telegramId int64) (user *User, err error) {
    pg := pgsql.GetClient()
    q := `select uuid, telegram_id, created_at from bloodpressure.public.user where telegram_id = $1`

    var row *sql.Row
    row = pg.QueryRow(q, telegramId)

    user = &User{}
    if err = row.Scan(&user.UUID, &user.TelegramId, &user.CreatedAt); err != nil {
        return nil, err
    }

    return user, nil
}
