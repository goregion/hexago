package database

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
)

type NamedStmt = sqlx.NamedStmt

type Client struct {
	*sqlx.DB
}

func NewClient(ctx context.Context, driverName string, dataSourceName string) (*Client, func(), error) {
	db, err := sqlx.Connect(driverName, dataSourceName)
	if err != nil {
		return nil, nil, err
	}
	client := &Client{
		DB: db,
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := client.PingContext(ctx); err != nil {
		db.Close()
		return nil, func() {}, err
	}

	return client,
		func() {
			db.Close()
		},
		nil
}
