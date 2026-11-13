package identity

import "time"

type User struct {
	Id          string       `db:"id"`
	Name        string       `db:"name"`
	Email       string       `db:"email"`
	Historicals []Historical `db:"historicals"`
	CreatedAt   time.Time    `db:"created_at"`
	UpdatedAt   time.Time    `db:"updated_at"`
}
