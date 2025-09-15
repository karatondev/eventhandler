package service

import (
	"context"
	"encoding/json"
	"eventhandler/internal/provider"
	"eventhandler/model"
	"eventhandler/model/entity"
	"eventhandler/util"
	"fmt"
	"time"
)

const (
	waaKeyPrefix = "waa:%s"
	qrKeyPrefix  = "qr:%s"
)

func (s *service) HandleEvent(ctx context.Context, data *model.QueueEvent) error {

	senderJID := util.ExtractJIDPrefix(data.SenderJID)

	switch data.EventType {
	case model.EventTypeQR:
		var qe model.QREventData
		if err := json.Unmarshal(data.Data, &qe); err != nil {
			return err
		}

		key := fmt.Sprintf(qrKeyPrefix, data.SenderJID)
		err := s.redis.Set(ctx, key, qe.Code, time.Duration(util.Configuration.Redis.QRSpan)*time.Second).Err()
		if err != nil {
			s.logger.Errorfctx(provider.AppLog, ctx, false, "Failed save QR event to redis: %v", err)
			return err
		}

		s.logger.Infofctx(provider.AppLog, ctx, "QR event for senderJID %s saved to redis", data.SenderJID)
		return nil

	case model.EventTypeOutboundMessage:

		account, err := s.GetAccountBySenderJID(ctx, senderJID)
		if err != nil {
			return err
		}

		var message model.OutboundMessageData
		if err := json.Unmarshal(data.Data, &message); err != nil {
			return err
		}

		messageData, err := json.Marshal(message.Message)
		if err != nil {
			return err
		}

		req := entity.CreateMessageOutboundRequest{
			EventID:     data.EventID,
			AccountID:   account.AccountID,
			MessageID:   message.MessageID,
			Recipient:   message.To,
			MessageType: string(message.MessageType),
			SentAt:      data.Timestamp,
			Data:        messageData,
		}

		if err := s.inboundOutbound.SaveMessageOutbound(ctx, &req); err != nil {
			s.logger.Errorfctx(provider.AppLog, ctx, false, "Failed save outbound event: %v", err)
			return err
		}

		s.logger.Infofctx(provider.AppLog, ctx, "Outbound event for senderJID %s saved", data.SenderJID)
		return nil

	case model.EventTypeInboundMessage:

		account, err := s.GetAccountBySenderJID(ctx, senderJID)
		if err != nil {
			return err
		}

		var message model.MessageEventData
		if err := json.Unmarshal(data.Data, &message); err != nil {
			return err
		}

		req := entity.CreateMessageInboundRequest{
			EventID:     data.EventID,
			AccountID:   account.AccountID,
			FromMe:      message.Metadata.FromMe,
			MessageID:   message.Metadata.MessageID,
			Sender:      message.Sender,
			MessageType: string(message.MessageType),
			ReceivedAt:  data.Timestamp,
		}

		switch message.MessageType {
		case model.MessageTypeText:
			// Create a clean message data without sender and metadata
			messageData := model.MessageData{
				Content: message.Content,
			}

			req.Data = messageData.MustToBytes()

		case model.MessageTypeImage:
			var imageMessage model.ImageMessageData
			if err := json.Unmarshal(data.Data, &imageMessage); err != nil {
				return err
			}

			messageData := model.MessageData{
				Content:  imageMessage.Content,
				Caption:  imageMessage.Caption,
				MimeType: imageMessage.MimeType,
				FileSize: imageMessage.FileSize,
				FileURL:  imageMessage.FileURL,
			}

			req.Data = messageData.MustToBytes()

		case model.MessageTypeAudio:
			var audioMessage model.AudioMessageData
			if err := json.Unmarshal(data.Data, &audioMessage); err != nil {
				return err
			}

			messageData := model.MessageData{
				Content:  audioMessage.Content,
				Caption:  audioMessage.Caption,
				MimeType: audioMessage.MimeType,
				FileSize: audioMessage.FileSize,
				FileURL:  audioMessage.FileURL,
				Duration: audioMessage.Duration,
			}

			req.Data = messageData.MustToBytes()

		case model.MessageTypeVideo:
			var videoMessage model.VideoMessageData
			if err := json.Unmarshal(data.Data, &videoMessage); err != nil {
				return err
			}

			messageData := model.MessageData{
				Content:  videoMessage.Content,
				Caption:  videoMessage.Caption,
				MimeType: videoMessage.MimeType,
				FileSize: videoMessage.FileSize,
				FileURL:  videoMessage.FileURL,
				Duration: videoMessage.Duration,
			}

			req.Data = messageData.MustToBytes()

		case model.MessageTypeDocument:
			var documentMessage model.DocumentMessageData
			if err := json.Unmarshal(data.Data, &documentMessage); err != nil {
				return err
			}

			messageData := model.MessageData{
				Content:  documentMessage.Content,
				Caption:  documentMessage.Caption,
				FileName: documentMessage.FileName,
				MimeType: documentMessage.MimeType,
				FileSize: documentMessage.FileSize,
				FileURL:  documentMessage.FileURL,
			}

			req.Data = messageData.MustToBytes()

		case model.MessageTypeLocation:
			var locationMessage model.LocationMessageData
			if err := json.Unmarshal(data.Data, &locationMessage); err != nil {
				return err
			}

			messageData := model.MessageData{
				Content:   locationMessage.Content,
				Caption:   locationMessage.Caption,
				Latitude:  locationMessage.Latitude,
				Longitude: locationMessage.Longitude,
				Name:      locationMessage.Name,
				Address:   locationMessage.Address,
			}

			req.Data = messageData.MustToBytes()

		case model.MessageTypeReaction:
			var reactionMessage model.ReactionMessageData
			if err := json.Unmarshal(data.Data, &reactionMessage); err != nil {
				return err
			}

			messageData := model.MessageData{
				Content:      reactionMessage.Content,
				Caption:      reactionMessage.Caption,
				Text:         reactionMessage.Text,
				TargetKey:    reactionMessage.TargetKey,
				TargetSender: reactionMessage.TargetSender,
			}

			req.Data = messageData.MustToBytes()

		case model.MessageTypeButton:
			var buttonMessage model.ButtonResponseMessageData
			if err := json.Unmarshal(data.Data, &buttonMessage); err != nil {
				return err
			}

			messageData := model.MessageData{
				Content:          buttonMessage.Content,
				Caption:          buttonMessage.Caption,
				SelectedButtonID: buttonMessage.SelectedButtonID,
				DisplayText:      buttonMessage.DisplayText,
			}

			req.Data = messageData.MustToBytes()

		case model.MessageTypeList:
			var listMessage model.ListResponseMessageData
			if err := json.Unmarshal(data.Data, &listMessage); err != nil {
				return err
			}

			messageData := model.MessageData{
				Content:       listMessage.Content,
				Caption:       listMessage.Caption,
				Title:         listMessage.Title,
				Description:   listMessage.Description,
				SelectedRowID: listMessage.SelectedRowID,
			}

			req.Data = messageData.MustToBytes()

		default:
			s.logger.Errorfctx(provider.AppLog, ctx, false, "Unsupported message type: %s", message.MessageType)
			return nil
		}

		if err := s.inboundOutbound.SaveMessageInbound(ctx, &req); err != nil {
			s.logger.Errorfctx(provider.AppLog, ctx, false, "Failed save inbound event: %v", err)
			return err
		}

		s.logger.Infofctx(provider.AppLog, ctx, "Inbound event for senderJID %s saved", data.SenderJID)
		return nil

	case model.EventTypePairSuccess:
		// Handle pair success event
		var pairEvent model.PairSuccessEventData
		if err := json.Unmarshal(data.Data, &pairEvent); err != nil {
			return err
		}

		accountID := util.ExtractJIDPrefix(pairEvent.AccountJID)
		req := entity.WhatsAppAccountReq{
			AccountID:     accountID,
			PhoneNumber:   &pairEvent.PhoneNumber,
			SenderJID:     &data.SenderJID,
			ConnectStatus: data.EventType,
			ConnectedAt:   &data.Timestamp,
		}

		if err := s.whatsappRepo.UpdatePairingSuccess(ctx, req); err != nil {
			s.logger.Errorfctx(provider.AppLog, ctx, false, "Failed to update pairing success: %v", err)
			return err
		}

		// Invalidate cache in Redis
		if err := s.redis.Del(ctx, fmt.Sprintf(waaKeyPrefix, data.SenderJID)).Err(); err != nil {
			s.logger.Errorfctx(provider.AppLog, ctx, false, "Failed to invalidate cache for sender JID %s: %v", data.SenderJID, err)
			return err
		}

		s.logger.Infofctx(provider.AppLog, ctx, "Pair success event for senderJID %s updated in database", data.SenderJID)
		return nil

	case model.EventTypeConnected:

		account, err := s.GetAccountBySenderJID(ctx, senderJID)
		if err != nil {
			return err
		}

		// Handle connected event
		var pairEvent model.ConnectionEventData
		if err := json.Unmarshal(data.Data, &pairEvent); err != nil {
			return err
		}

		req := entity.WhatsAppAccountReq{
			AccountID:     account.AccountID,
			ConnectStatus: data.EventType,
			ConnectedAt:   &data.Timestamp,
		}

		if err := s.whatsappRepo.UpdateConnected(ctx, req); err != nil {
			s.logger.Errorfctx(provider.AppLog, ctx, false, "Failed to update connected status: %v", err)
			return err
		}

		s.logger.Infofctx(provider.AppLog, ctx, "Connected event for senderJID %s updated in database", data.SenderJID)
		return nil

	case model.EventTypeDisconnected:

		account, err := s.GetAccountBySenderJID(ctx, senderJID)
		if err != nil {
			return err
		}

		// Handle disconnected event
		var pairEvent model.ConnectionEventData
		if err := json.Unmarshal(data.Data, &pairEvent); err != nil {
			return err
		}

		req := entity.WhatsAppAccountReq{
			AccountID:      account.AccountID,
			ConnectStatus:  data.EventType,
			DisconnectedAt: &data.Timestamp,
		}

		if err := s.whatsappRepo.UpdateDisconnected(ctx, req); err != nil {
			s.logger.Errorfctx(provider.AppLog, ctx, false, "Failed to update disconnected status: %v", err)
			return err
		}

		s.logger.Infofctx(provider.AppLog, ctx, "Disconnected event for senderJID %s updated in database", data.SenderJID)
		return nil

	case model.EventTypeLoggedOut:

		account, err := s.GetAccountBySenderJID(ctx, senderJID)
		if err != nil {
			return err
		}

		// Handle logged out event - same as disconnected
		var pairEvent model.ConnectionEventData
		if err := json.Unmarshal(data.Data, &pairEvent); err != nil {
			return err
		}

		req := entity.WhatsAppAccountReq{
			AccountID:      account.AccountID,
			ConnectStatus:  data.EventType,
			DisconnectedAt: &data.Timestamp,
		}

		if err := s.whatsappRepo.UpdateDisconnected(ctx, req); err != nil {
			s.logger.Errorfctx(provider.AppLog, ctx, false, "Failed to update logged out status: %v", err)
			return err
		}

		s.logger.Infofctx(provider.AppLog, ctx, "Logged out event for senderJID %s updated in database", data.SenderJID)
		return nil

	case model.EventTypeReceipt:

		account, err := s.GetAccountBySenderJID(ctx, senderJID)
		if err != nil {
			return err
		}

		var receipt model.ReceiptEventData
		if err := json.Unmarshal(data.Data, &receipt); err != nil {
			return err
		}

		for _, messageID := range receipt.MessageIDs {
			req := entity.CreateMessageReceiptRequest{
				ReceiptID: data.EventID,
				AccountID: account.AccountID,
				MessageID: messageID,
				Timestamp: time.Unix(receipt.Timestamp, 0),
				Status:    receipt.Type,
			}

			if err := s.inboundOutbound.SaveMessageReceipt(ctx, &req); err != nil {
				s.logger.Errorfctx(provider.AppLog, ctx, false, "Failed save message receipt: %v", err)
				return err
			}
		}

		s.logger.Infofctx(provider.AppLog, ctx, "Receipt event for accountID %s saved (%d receipts)", account.AccountID, len(receipt.MessageIDs))
		return nil

	case model.EventTypePresence, model.EventTypeCallOffer, model.EventTypeMediaRetryError:

		account, err := s.GetAccountBySenderJID(ctx, senderJID)
		if err != nil {
			return err
		}

		req := entity.CreateAccountEventRequest{
			EventID:   data.EventID,
			AccountID: account.AccountID,
			EventType: string(data.EventType),
			Timestamp: data.Timestamp,
			Data:      data.Data,
		}

		if err := s.inboundOutbound.SaveEvent(ctx, &req); err != nil {
			s.logger.Errorfctx(provider.AppLog, ctx, false, "Failed save event: %v", err)
			return err
		}
		s.logger.Infofctx(provider.AppLog, ctx, "Event %s for accountID %s saved", data.EventType, account.AccountID)
		return nil

	default:
		return fmt.Errorf("unknown event type: %s", data.EventType)

	}
}

func (s *service) GetAccountBySenderJID(ctx context.Context, senderJID string) (*entity.WhatsAppAccount, error) {
	// Try to get from Redis first
	key := fmt.Sprintf(waaKeyPrefix, senderJID)
	cachedData, err := s.redis.Get(ctx, key).Result()
	if err == nil {
		// Data found in Redis, unmarshal and return
		var account entity.WhatsAppAccount
		if err := json.Unmarshal([]byte(cachedData), &account); err == nil {
			s.logger.Infofctx(provider.AppLog, ctx, "WhatsApp account found in Redis for JID: %s", senderJID)
			return &account, nil
		}
		// If unmarshal fails, continue to database fetch
		s.logger.Errorfctx(provider.AppLog, ctx, false, "Failed to unmarshal cached account data: %v", err)
	}

	// Data not found in Redis or unmarshal failed, fetch from database
	account, err := s.whatsappRepo.GetAccountBySenderJID(ctx, senderJID)
	if err != nil {
		s.logger.Errorfctx(provider.AppLog, ctx, false, "Failed to get account from database: %v", err)
		return nil, err
	}

	// Cache the result in Redis for future use (cache for 1 hour)
	accountData, err := json.Marshal(account)
	if err == nil {
		s.redis.Set(ctx, key, accountData, time.Second).Err()
		s.logger.Infofctx(provider.AppLog, ctx, "WhatsApp account cached in Redis for JID: %s", senderJID)
	}

	s.logger.Infofctx(provider.AppLog, ctx, "WhatsApp account found in database for JID: %s", senderJID)
	return account, nil
}
