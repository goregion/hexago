package adapter_mysql

import (
	"context"

	"github.com/goregion/hexago/internal/entity"
	sqlgen_db "github.com/goregion/hexago/pkg/sqlgen-db"
)

type OHLCRepository struct {
	databaseClient *sqlgen_db.Client
}

func NewOHLCRepository(databaseClient *sqlgen_db.Client, timeframeName string) *OHLCRepository {
	return &OHLCRepository{
		databaseClient: databaseClient,
	}
}

// StoreOHLC inserts a new OHLC record into the database
func (p *OHLCRepository) StoreOHLC(ctx context.Context, ohlc *entity.OHLC) error {
	executor := sqlgen_db.GetTransaction(ctx, p.databaseClient)
	if _, err := executor.Insert("ohlc").Rows(ohlc).Executor().ExecContext(ctx); err != nil {
		return err
	}
	return nil
}
