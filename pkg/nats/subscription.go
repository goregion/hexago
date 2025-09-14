package natsmq

import (
	"context"
	"errors"
	"time"

	"github.com/nats-io/nats.go"
)

type subscription struct {
	*nats.Subscription
}

func (s *subscription) ReadNext(ctx context.Context, timeout time.Duration) (string, []byte, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	msg, err := s.NextMsgWithContext(ctx)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			err = nil
		}
		return "", nil, err
	}
	return msg.Subject, msg.Data, nil
}
