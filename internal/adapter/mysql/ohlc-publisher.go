package adapter_mysql

import (
	"context"

	"github.com/goregion/hexago/internal/entity"
	"github.com/goregion/hexago/pkg/database"
	"github.com/goregion/must"
)

type OHLCPublisher struct {
	databaseClient *database.Client
	insertStmt     *database.NamedStmt
}

func NewOHLCPublisher(ctx context.Context, databaseClient *database.Client, timeframeName string) *OHLCPublisher {
	return &OHLCPublisher{
		databaseClient: databaseClient,
		insertStmt: must.Return(
			databaseClient.PrepareNamedContext(ctx,
				"INSERT INTO "+makeOHLCTableName(timeframeName)+" (symbol, open, high, low, close, close_time_ms, timeframe) VALUES (:symbol, :open, :high, :low, :close, :close_time_ms, :timeframe)",
			),
		),
	}
}

// PublishOHLC publishes the given OHLC to the appropriate MySQL table
func (p *OHLCPublisher) PublishOHLC(ctx context.Context, ohlc *entity.OHLC) error {
	if _, err := p.insertStmt.ExecContext(ctx, ohlc); err != nil {
		return err
	}
	return nil
}
