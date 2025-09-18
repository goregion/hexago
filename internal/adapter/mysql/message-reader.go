package adapter_mysql

import (
	"context"

	"github.com/goregion/hexago/internal/entity"
	"github.com/goregion/hexago/pkg/database"
	"github.com/goregion/must"
)

type MessageReader struct {
	databaseClient *database.Client

	selectStmt *database.NamedStmt
}

func NewMessageReader(ctx context.Context, databaseClient *database.Client) *MessageReader {
	return &MessageReader{
		databaseClient: databaseClient,
		selectStmt: must.Return(
			databaseClient.PrepareNamedContext(ctx,
				"SELECT `value` FROM "+messagesTable+" WHERE `id` = :id",
			),
		),
	}
}

func (p *MessageReader) ReadMessage(ctx context.Context, id int) (*entity.Message, error) {
	var message entity.Message
	if err := p.selectStmt.GetContext(ctx,
		&message,
		map[string]any{
			"id": id,
		},
	); err != nil {
		return nil, err
	}
	return &message, nil
}
