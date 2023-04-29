package storage

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/meshenka/hunt/generated"
)

type DBTX interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

// DB is a connection pool that we can create transactions from.
type DB interface {
	DBTX
	BeginTx(context.Context, *sql.TxOptions) (Tx, error)
}

// Tx is a transaction as an inteface so we can mock it.
type Tx interface {
	DBTX
	Commit() error
	Rollback() error
}

type DBAdapter struct{ *sql.DB } // DBAdapter is an implementation of DB based off *sql.DB.

// BeginTx implements DB for DBAdapter.
func (db DBAdapter) BeginTx(ctx context.Context, opts *sql.TxOptions) (Tx, error) {
	return db.DB.BeginTx(ctx, opts)
}

type Storage struct {
	db DB
}

func New(db DB) Storage {
	return Storage{db: db}
}

type OpportunityAggregate struct {
	Author      generated.Account
	Opportunity generated.Opportunity
	Notes       []generated.OpportunityNote
	Company     generated.Company
}

func (s Storage) LoadOpportuniry(ctx context.Context, authorID uuid.UUID, opportunityID uuid.UUID) (OpportunityAggregate, error) {
	q := generated.New(s.db)

	author, err := q.AccountByID(ctx, authorID)
	if err != nil {
		return OpportunityAggregate{}, fmt.Errorf("cannot get author: %w", err)
	}

	if !author.Enabled {
		return OpportunityAggregate{}, fmt.Errorf("author is disabled")
	}

	if !author.Deleted {
		return OpportunityAggregate{}, fmt.Errorf("author account is deleted")
	}

	opportunity, err := q.OpportunityByID(ctx, generated.OpportunityByIDParams{
		AuthorID: authorID,
		ID:       opportunityID,
	})
	if err != nil {
		return OpportunityAggregate{}, fmt.Errorf("cannot get opportunity: %w", err)
	}

	company, err := q.CompanyByID(ctx, generated.CompanyByIDParams{
		AuthorID: authorID,
		ID:       opportunity.CompanyID,
	})
	if err != nil {
		return OpportunityAggregate{}, fmt.Errorf("cannot get company: %w", err)
	}

	notes, err := q.OpportunityNotes(ctx, generated.OpportunityNotesParams{
		AuthorID:      authorID,
		OpportunityID: opportunityID,
	})
	if err != nil {
		return OpportunityAggregate{}, fmt.Errorf("cannot get notes: %w", err)
	}

	return OpportunityAggregate{
		Author:      author,
		Company:     company,
		Opportunity: opportunity,
		Notes:       notes,
	}, nil

}
func (s Storage) Companies(ctx context.Context, authorID uuid.UUID) ([]generated.Company, error) {
	q := generated.New(s.db)
	return q.Companies(ctx, authorID)
}
