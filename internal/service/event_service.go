package service

import (
	"context"
	"encoding/json"
	"eventhandler/model"
	"eventhandler/model/entity"
	"eventhandler/util"
	"fmt"
	"time"

	sharedmodel "zaplio/shared/model"
	"zaplio/shared/pkg/logger"

	"github.com/google/uuid"
)

const (
	waaKeyPrefix = "waa:%s"
	qrKeyPrefix  = "qr:%s"
)

func (s *service) HandleEvent(ctx context.Context, data *sharedmodel.QueueEvent) error {

	senderJID := util.ExtractJIDPrefix(data.SenderJID)
	s.logger.Infofctx(logger.AppLog, ctx, "Processing event: event_id=%s, event_type=%s, sender_jid=%s", data.EventID, data.EventType, data.SenderJID)

	switch data.EventType {
	case sharedmodel.EventTypeQR:
		var qe sharedmodel.QREventData
		if err := json.Unmarshal(data.Data, &qe); err != nil {
			return err
		}

		key := fmt.Sprintf(qrKeyPrefix, senderJID)
		err := s.redis.Set(ctx, key, qe.Code, time.Duration(util.Configuration.Redis.QRSpan)*time.Second).Err()
		if err != nil {
			s.logger.Errorfctx(logger.AppLog, ctx, false, "Failed save QR event to redis: %v", err)
			return err
		}

		s.logger.Infofctx(logger.AppLog, ctx, "QR event for senderJID %s saved to redis", senderJID)
		return nil

	case sharedmodel.EventTypeOutboundMessage:

		account, err := s.GetAccountBySenderJID(ctx, senderJID)
		if err != nil {
			return err
		}

		var message sharedmodel.OutboundMessageData
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

		s.logger.Infofctx(logger.AppLog, ctx, "Saving outbound message: event_id=%s, account_id=%s, message_id=%s, recipient=%s, message_type=%s", req.EventID, req.AccountID, req.MessageID, req.Recipient, req.MessageType)
		if err := s.inboundOutbound.SaveMessageOutbound(ctx, &req); err != nil {
			s.logger.Errorfctx(logger.AppLog, ctx, false, "Failed save outbound event: %v", err)
			return err
		}

		s.logger.Infofctx(logger.AppLog, ctx, "Outbound event for senderJID %s saved", data.SenderJID)
		return nil

	case sharedmodel.EventTypeInboundMessage:

		account, err := s.GetAccountBySenderJID(ctx, senderJID)
		if err != nil {
			return err
		}

		var message sharedmodel.MessageEventData
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
		case sharedmodel.MessageTypeText:
			// Create a clean message data without sender and metadata
			messageData := model.MessageData{
				Content: message.Content,
			}

			req.Data = messageData.MustToBytes()

		case sharedmodel.MessageTypeImage:
			var imageMessage sharedmodel.ImageMessageData
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

		case sharedmodel.MessageTypeAudio:
			var audioMessage sharedmodel.AudioMessageData
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

		case sharedmodel.MessageTypeVideo:
			var videoMessage sharedmodel.VideoMessageData
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

		case sharedmodel.MessageTypeDocument:
			var documentMessage sharedmodel.DocumentMessageData
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

		case sharedmodel.MessageTypeLocation:
			var locationMessage sharedmodel.LocationMessageData
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

		case sharedmodel.MessageTypeReaction:
			var reactionMessage sharedmodel.ReactionMessageData
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

		case sharedmodel.MessageTypeButton:
			var buttonMessage sharedmodel.ButtonResponseMessageData
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

		case sharedmodel.MessageTypeList:
			var listMessage sharedmodel.ListResponseMessageData
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
			s.logger.Errorfctx(logger.AppLog, ctx, false, "Unsupported message type: %s", message.MessageType)
			return nil
		}

		s.logger.Infofctx(logger.AppLog, ctx, "Saving inbound message: event_id=%s, account_id=%s, message_id=%s, sender=%s, message_type=%s", req.EventID, req.AccountID, req.MessageID, req.Sender, req.MessageType)
		if err := s.inboundOutbound.SaveMessageInbound(ctx, &req); err != nil {
			s.logger.Errorfctx(logger.AppLog, ctx, false, "Failed save inbound event: %v", err)
			return err
		}

		s.logger.Infofctx(logger.AppLog, ctx, "Inbound event for senderJID %s saved", data.SenderJID)
		return nil

	case sharedmodel.EventTypePairSuccess:
		// Handle pair success event
		var pairEvent sharedmodel.PairSuccessEventData
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
			s.logger.Errorfctx(logger.AppLog, ctx, false, "Failed to update pairing success: %v", err)
			return err
		}

		// Invalidate cache in Redis
		if err := s.redis.Del(ctx, fmt.Sprintf(waaKeyPrefix, data.SenderJID)).Err(); err != nil {
			s.logger.Errorfctx(logger.AppLog, ctx, false, "Failed to invalidate cache for sender JID %s: %v", data.SenderJID, err)
			return err
		}

		s.logger.Infofctx(logger.AppLog, ctx, "Pair success event for senderJID %s updated in database", data.SenderJID)
		return nil

	case sharedmodel.EventTypeConnected:

		account, err := s.GetAccountBySenderJID(ctx, senderJID)
		if err != nil {
			return err
		}

		// Handle connected event
		var pairEvent sharedmodel.ConnectionEventData
		if err := json.Unmarshal(data.Data, &pairEvent); err != nil {
			return err
		}

		req := entity.WhatsAppAccountReq{
			AccountID:     account.AccountID,
			ConnectStatus: data.EventType,
			ConnectedAt:   &data.Timestamp,
		}

		if err := s.whatsappRepo.UpdateConnected(ctx, req); err != nil {
			s.logger.Errorfctx(logger.AppLog, ctx, false, "Failed to update connected status: %v", err)
			return err
		}

		s.logger.Infofctx(logger.AppLog, ctx, "Connected event for senderJID %s updated in database", data.SenderJID)
		return nil

	case sharedmodel.EventTypeDisconnected:

		account, err := s.GetAccountBySenderJID(ctx, senderJID)
		if err != nil {
			return err
		}

		// Handle disconnected event
		var pairEvent sharedmodel.ConnectionEventData
		if err := json.Unmarshal(data.Data, &pairEvent); err != nil {
			return err
		}

		req := entity.WhatsAppAccountReq{
			AccountID:      account.AccountID,
			ConnectStatus:  data.EventType,
			DisconnectedAt: &data.Timestamp,
		}

		if err := s.whatsappRepo.UpdateDisconnected(ctx, req); err != nil {
			s.logger.Errorfctx(logger.AppLog, ctx, false, "Failed to update disconnected status: %v", err)
			return err
		}

		s.logger.Infofctx(logger.AppLog, ctx, "Disconnected event for senderJID %s updated in database", data.SenderJID)
		return nil

	case sharedmodel.EventTypeLoggedOut:

		account, err := s.GetAccountBySenderJID(ctx, senderJID)
		if err != nil {
			return err
		}

		// Handle logged out event - same as disconnected
		var pairEvent sharedmodel.ConnectionEventData
		if err := json.Unmarshal(data.Data, &pairEvent); err != nil {
			return err
		}

		req := entity.WhatsAppAccountReq{
			AccountID:      account.AccountID,
			ConnectStatus:  data.EventType,
			DisconnectedAt: &data.Timestamp,
		}

		if err := s.whatsappRepo.UpdateDisconnected(ctx, req); err != nil {
			s.logger.Errorfctx(logger.AppLog, ctx, false, "Failed to update logged out status: %v", err)
			return err
		}

		s.logger.Infofctx(logger.AppLog, ctx, "Logged out event for senderJID %s updated in database", data.SenderJID)
		return nil

	case sharedmodel.EventTypeReceipt:

		account, err := s.GetAccountBySenderJID(ctx, senderJID)
		if err != nil {
			return err
		}

		var receipt sharedmodel.ReceiptEventData
		if err := json.Unmarshal(data.Data, &receipt); err != nil {
			return err
		}

		for _, messageID := range receipt.MessageIDs {
			// Satu event receipt bisa memuat banyak message_id; receipt_id harus unik
			// per (event, message) agar tidak bentrok PK. Deterministik supaya redelivery
			// menghasilkan id yang sama (idempotent dengan ON CONFLICT DO NOTHING di repo).
			receiptID := uuid.NewSHA1(uuid.NameSpaceOID, []byte(data.EventID+":"+messageID)).String()

			req := entity.CreateMessageReceiptRequest{
				ReceiptID: receiptID,
				AccountID: account.AccountID,
				MessageID: messageID,
				Timestamp: time.Unix(receipt.Timestamp, 0),
				Status:    receipt.Type,
			}

			if err := s.inboundOutbound.SaveMessageReceipt(ctx, &req); err != nil {
				s.logger.Errorfctx(logger.AppLog, ctx, false, "Failed save message receipt: %v", err)
				return err
			}
		}

		s.logger.Infofctx(logger.AppLog, ctx, "Receipt event for accountID %s saved (%d receipts)", account.AccountID, len(receipt.MessageIDs))
		return nil

	case sharedmodel.EventTypePresence, sharedmodel.EventTypeCallOffer, sharedmodel.EventTypeMediaRetryError:

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
			s.logger.Errorfctx(logger.AppLog, ctx, false, "Failed save event: %v", err)
			return err
		}
		s.logger.Infofctx(logger.AppLog, ctx, "Event %s for accountID %s saved", data.EventType, account.AccountID)
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
			s.logger.Infofctx(logger.AppLog, ctx, "WhatsApp account found in Redis for JID: %s", senderJID)
			return &account, nil
		}
		// If unmarshal fails, continue to database fetch
		s.logger.Errorfctx(logger.AppLog, ctx, false, "Failed to unmarshal cached account data: %v", err)
	}

	// Data not found in Redis or unmarshal failed, fetch from database
	account, err := s.whatsappRepo.GetAccountBySenderJID(ctx, senderJID)
	if err != nil {
		s.logger.Errorfctx(logger.AppLog, ctx, false, "Failed to get account from database: %v", err)
		return nil, err
	}

	accountData, err := json.Marshal(account)
	if err == nil {
		s.redis.Set(ctx, key, accountData, time.Hour).Err()
		s.logger.Infofctx(logger.AppLog, ctx, "WhatsApp account cached in Redis for JID: %s", senderJID)
	}

	s.logger.Infofctx(logger.AppLog, ctx, "WhatsApp account found in database for JID: %s", senderJID)
	return account, nil
}
