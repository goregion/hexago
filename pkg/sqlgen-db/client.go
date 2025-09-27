package sqlgen_db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/doug-martin/goqu/v9"
)

type Transaction interface {
	From(table ...interface{}) *goqu.SelectDataset
	Update(table interface{}) *goqu.UpdateDataset
	Insert(table interface{}) *goqu.InsertDataset
}

type Client struct {
	*goqu.Database
}

func NewClient(ctx context.Context, driverName string, dataSourceName string) (*Client, func(), error) {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, nil, err
	}
	var goquDB = goqu.New(driverName, db)

	client := &Client{
		Database: goquDB,
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, func() {}, err
	}

	return client,
		func() {
			db.Close()
		},
		nil
}

type CommitTxFunc = func()
type RollbackTxFunc = func()

type txContextKeyType string

const txContextKey txContextKeyType = "sql-db-tx"

func (db *Client) WithTx(ctx context.Context) (context.Context, CommitTxFunc, RollbackTxFunc, error) {
	// Get underlying *sql.DB from goqu.Database
	sqlDB := db.Database.Db

	sqlTx, err := sqlDB.BeginTx(ctx,
		&sql.TxOptions{
			Isolation: sql.LevelReadCommitted,
			ReadOnly:  false,
		},
	)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Create goqu TxDatabase for working with transaction
	goquTx := goqu.NewTx(db.Database.Dialect(), sqlTx)
	txCtx := context.WithValue(ctx, txContextKey, goquTx)

	var finished bool

	commitFunc := func() error {
		if finished {
			return fmt.Errorf("transaction already finished")
		}
		finished = true
		return sqlTx.Commit()
	}

	rollbackFunc := func() error {
		if finished {
			return fmt.Errorf("transaction already finished")
		}
		finished = true
		return sqlTx.Rollback()
	}

	return txCtx,
		func() {
			if err := commitFunc(); err != nil {
				// In case of commit error, try to rollback
				rollbackFunc()
				panic(fmt.Errorf("failed to commit transaction: %w", err))
			}
		},
		func() {
			if err := rollbackFunc(); err != nil {
				panic(fmt.Errorf("failed to rollback transaction: %w", err))
			}
		},
		nil
}

func GetTransaction(ctx context.Context, client *Client) Transaction {
	if tx, ok := ctx.Value(txContextKey).(*goqu.TxDatabase); ok {
		return tx
	}
	return client
}
