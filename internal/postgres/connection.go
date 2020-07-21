package postgresrepo

import (
	"context"
	"database/sql"
	"fmt"

	// postgres driver
	_ "github.com/lib/pq"
)

// Repository with users' credentials
type Repository struct {
	db *sql.DB
}

// Configuration is settings for postgress connection
type Configuration struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

// New returns new instance of repository
func New(ctx context.Context, conf *Configuration) (*Repository, error) {
	host := "localhost"
	var user, dbname string = "postgres", "postgres"
	password := ""
	port := 5432

	// TODO: parse URI

	if conf != nil {
		if conf.Host != "" {
			host = conf.Host
		}
		if conf.User != "" {
			user = conf.User
		}
		if conf.DBName != "" {
			dbname = conf.DBName
		}
		if conf.Port != 0 {
			port = conf.Port
		}

		password = conf.Password

	}

	db, err := sql.Open("postgres",
		fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			host,
			port,
			user,
			password,
			dbname,
		))
	if err != nil {
		return nil, err
	}

	_, err = db.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS test_transactions (
			id					VARCHAR,
			user_id			VARCHAR,
			created_at	TIMESTAMPTZ
		);`)
	if err != nil {
		return &Repository{
			db: db,
		}, err
	}

	return &Repository{
		db: db,
	}, nil
}
