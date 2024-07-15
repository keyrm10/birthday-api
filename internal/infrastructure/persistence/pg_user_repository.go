package persistence

import (
	"database/sql"
	"time"

	_ "github.com/lib/pq"

	"github.com/keyrm10/birthday-api/internal/domain/user"
)

type PgUserRepository struct {
	db *sql.DB
}

func NewPgUserRepository(db *sql.DB) *PgUserRepository {
	return &PgUserRepository{db: db}
}

func (r *PgUserRepository) Save(u user.User) error {
	_, err := r.db.Exec("INSERT INTO users (username, date_of_birth) VALUES ($1, $2) ON CONFLICT (username) DO UPDATE SET date_of_birth = $2",
		string(u.Username), u.DateOfBirth.Time.Format(time.DateOnly))
	return err
}

func (r *PgUserRepository) FindByUsername(username user.Username) (user.User, error) {
	var u user.User
	var dateOfBirth time.Time
	err := r.db.QueryRow("SELECT username, date_of_birth FROM users WHERE username = $1", string(username)).Scan(&u.Username, &dateOfBirth)
	if err != nil {
		return user.User{}, err
	}
	u.DateOfBirth = user.DateOfBirth{Time: dateOfBirth}
	return u, nil
}

func NewPostgresDB(dataSourceName string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
