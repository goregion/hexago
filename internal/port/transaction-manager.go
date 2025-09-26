package port

import "context"

type CommitTxFunc = func()
type RollbackTxFunc = func()

type TransactionManager interface {
	WithTx(ctx context.Context) (context.Context, CommitTxFunc, RollbackTxFunc, error)
}
