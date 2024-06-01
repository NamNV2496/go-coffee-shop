package domain

import (
	"time"

	"github.com/doug-martin/goqu/v9"
)

var (
	TabNameUser = goqu.T("user")
	ColId       = "id"
	ColUserId   = "user_id"
)

type User struct {
	Id          int       `db:"id" goqu:"omitnil"`
	UserId      string    `db:"user_id" goqu:"omitnil"`
	Password    string    `db:"password" goqu:"omitnil"`
	Name        string    `db:"name" goqu:"omitnil"`
	Age         int       `db:"age" goqu:"omitnil"`
	Position    string    `db:"position" goqu:"omitnil"`
	Email       string    `db:"email" goqu:"omitnil"`
	IsActive    bool      `db:"is_active" goqu:"omitnil"`
	Role        string    `db:"role" goqu:"omitnil"`
	CreatedDate time.Time `db:"created_date" goqu:"omitnil"`
}
