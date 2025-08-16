package service

import (
	"context"
	"encoding/json"
	"eventhandler/internal/provider"
	"eventhandler/model"
	"eventhandler/util"
	"fmt"
	"time"
)

func (s *service) HandleEvent(ctx context.Context, data *model.QueueEvent) error {
	switch data.EventType {
	case model.EventTypeQR:
		var qe model.QREventData
		if err := json.Unmarshal(data.Data, &qe); err != nil {
			return err
		}

		key := "qr:" + data.DeviceID
		err := s.redis.Set(ctx, key, qe.Code, time.Duration(util.Configuration.Redis.QRSpan)*time.Second).Err()
		if err != nil {
			s.logger.Errorfctx(provider.AppLog, ctx, false, "Failed save QR event to redis: %v", err)
			return err
		}

		s.logger.Infofctx(provider.AppLog, ctx, "QR event for device %s saved to redis", data.DeviceID)
		return nil

	case model.EventTypeMessage:
		if err := s.repo.SaveInbound(ctx, data); err != nil {
			s.logger.Errorfctx(provider.AppLog, ctx, false, "Failed save inbound event: %v", err)
			return err
		}
		s.logger.Infofctx(provider.AppLog, ctx, "Inbound event for device %s saved", data.DeviceID)
		return nil

	case model.EventTypeConnected, model.EventTypeDisconnected, model.EventTypeLoggedOut, model.EventTypePairSuccess, model.EventTypeReceipt, model.EventTypePresence, model.EventTypeCallOffer, model.EventTypeMediaRetryError:
		if err := s.repo.SaveEvent(ctx, data); err != nil {
			s.logger.Errorfctx(provider.AppLog, ctx, false, "Failed save event: %v", err)
			return err
		}
		s.logger.Infofctx(provider.AppLog, ctx, "Event %s for device %s saved", data.EventType, data.DeviceID)
		return nil

	default:
		return fmt.Errorf("unknown event type: %s", data.EventType)

	}
}
