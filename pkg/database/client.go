package database

import (
	"context"

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
	return client,
		func() {
			db.Close()
		},
		client.PingContext(ctx)
}
