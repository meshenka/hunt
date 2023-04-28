package storage

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/meshenka/hunt/generated"
)

type DBTX interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}
type Storage struct {
   db DBTX
}

func New(db DBTX) Storage {
  return Storage{db: db}
}

func (s Storage) Companies(ctx context.Context, authorID uuid.UUID) ([]generated.Company, error) {
   q := generated.New(s.db)
   return q.Companies(ctx, authorID)
}
