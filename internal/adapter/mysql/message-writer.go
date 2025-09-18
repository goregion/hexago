package adapter_mysql

import (
	"context"

	"github.com/goregion/hexago/internal/entity"
	"github.com/goregion/hexago/pkg/database"
	"github.com/goregion/must"
)

type MessagePublisher struct {
	databaseClient *database.Client

	insertStmt *database.NamedStmt
}

func NewMessagePublisher(ctx context.Context, databaseClient *database.Client) *MessagePublisher {
	return &MessagePublisher{
		databaseClient: databaseClient,
		insertStmt: must.Return(
			databaseClient.PrepareNamedContext(ctx,
				"INSERT INTO "+messagesTable+" (`value`) VALUES (:value)",
			),
		),
	}
}

func (p *MessagePublisher) PublishMessage(ctx context.Context, message *entity.Message) error {
	if _, err := p.insertStmt.ExecContext(ctx, message); err != nil {
		return err
	}
	return nil
}
