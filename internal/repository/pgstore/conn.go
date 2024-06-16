package pgstore

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/lib/pq"
)

const (
	Driver_Name = "postgres"
)

type DBTX interface {
	QueryContext(context.Context, string, ...any) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...any) *sql.Row
	ExecContext(context.Context, string, ...any) (sql.Result, error)
}

type Queries struct {
	db DBTX
}

func NewQueries(db DBTX) *Queries {
	return &Queries{
		db: db,
	}
}

type PGStore struct {
	*Queries

	db *sql.DB
}

type config interface {
	MustLoadString(string) string
	LoadInt(string, int) int
}

func New(config config) (*PGStore, error) {
	uri := config.MustLoadString("POSTGRES_URI")

	db, err := sql.Open(Driver_Name, uri)
	if err != nil {
		return nil, err
	}

	if err := db.PingContext(context.Background()); err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(config.LoadInt("DB_MAX_OPEN_CONNS", 25))
	db.SetMaxIdleConns(config.LoadInt("DB_MAX_IDLE_CONNS", 25))
	db.SetConnMaxIdleTime(time.Duration(config.LoadInt("CONN_MAX_IDLE_TIME", 15)))

	return &PGStore{Queries: NewQueries(db), db: db}, nil
}

func (store *PGStore) Close() error {
	return store.db.Close()
}
