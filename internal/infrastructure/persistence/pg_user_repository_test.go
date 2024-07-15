package persistence

import (
	"database/sql"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/assert"

	"github.com/keyrm10/birthday-api/internal/domain/user"
)

var db *sql.DB

func createTable(db *sql.DB) error {
	_, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS users (
            username VARCHAR(255) PRIMARY KEY,
            date_of_birth DATE NOT NULL
        );
		CREATE INDEX IF NOT EXISTS idx_users_date_of_birth ON users(date_of_birth);
    `)
	return err
}

func clearTable(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM users")
	return err
}

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("could not construct pool: %s", err)
	}

	err = pool.Client.Ping()
	if err != nil {
		log.Fatalf("could not connect to Docker: %s", err)
	}

	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "15",
		Env: []string{
			"POSTGRES_PASSWORD=secret",
			"POSTGRES_USER=user",
			"POSTGRES_DB=testdb",
			"listen_addresses = '*'",
		},
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		log.Fatalf("could not start resource: %s", err)
	}

	hostAndPort := resource.GetHostPort("5432/tcp")
	databaseUrl := fmt.Sprintf("postgres://user:secret@%s/testdb?sslmode=disable", hostAndPort)

	log.Println("Connecting to database on url: ", databaseUrl)
	if err := resource.Expire(120); err != nil {
		log.Fatalf("could not set expiration on resource: %s", err)
	}

	if err = pool.Retry(func() error {
		db, err = sql.Open("postgres", databaseUrl)
		if err != nil {
			return err
		}
		if err = db.Ping(); err != nil {
			return err
		}
		return createTable(db)
	}); err != nil {
		log.Fatalf("could not connect to Docker or create table: %s", err)
	}

	defer func() {
		if err := pool.Purge(resource); err != nil {
			log.Fatalf("could not purge resource: %s", err)
		}
	}()

	m.Run()
}

func TestPgUserRepository(t *testing.T) {
	repo := NewPgUserRepository(db)

	t.Run("SaveAndFindByUsername", func(t *testing.T) {
		err := clearTable(db)
		assert.NoError(t, err)

		username, _ := user.NewUsername("john")
		dob, _ := user.NewDateOfBirth("1990-01-01")
		u := user.User{Username: username, DateOfBirth: dob}

		err = repo.Save(u)
		assert.NoError(t, err)

		found, err := repo.FindByUsername(username)
		assert.NoError(t, err)
		assert.Equal(t, u.Username, found.Username)
		assert.Equal(t, u.DateOfBirth.Format(time.DateOnly), found.DateOfBirth.Format(time.DateOnly))
	})

	t.Run("FindByUsername_NotFound", func(t *testing.T) {
		err := clearTable(db)
		assert.NoError(t, err)

		username, _ := user.NewUsername("nonexistent")
		_, err = repo.FindByUsername(username)
		assert.Error(t, err)
	})
}
