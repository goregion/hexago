package adapter_nats

import (
	"context"

	"feeder/internal/entity"
	"feeder/internal/port"
	natsmq "feeder/pkg/nats"

	"google.golang.org/protobuf/encoding/protojson"
)

type TradeRecordSubscriber struct {
	subscription       natsmq.SubscriptionInterface
	tradeRecordHandler port.TradeRecordHandler
}

func NewTradeRecordSubscriber(
	subscription natsmq.SubscriptionInterface,
	tradeRecordHandler port.TradeRecordHandler,
) *TradeRecordSubscriber {
	return &TradeRecordSubscriber{
		subscription:       subscription,
		tradeRecordHandler: tradeRecordHandler,
	}
}

func (s *TradeRecordSubscriber) ReadNext(ctx context.Context) error {
	_, data, err := s.subscription.ReadNext(ctx, natsmq.DefaultTimeout)
	if err != nil {
		return err
	}
	if data == nil {
		return nil
	}

	var dataProto = &entity.TradeRecord{}
	if err := protojson.Unmarshal(data, dataProto); err != nil {
		return err
	}

	if err := s.tradeRecordHandler.HandleTradeRecord(ctx, dataProto); err != nil {
		return err
	}

	return nil
}
