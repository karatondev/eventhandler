package repository

import (
	"context"
	"eventhandler/model"
	"eventhandler/util"

	"github.com/jackc/pgx/v5"
)

// SaveEvent saves an inbound event to the database

const (
	insertEventQuery = `
       INSERT INTO "whatsapp-web".events 
	       (event_id, device_id, event_type, timestamp, data, created_at, updated_at) 
       VALUES (@event_id, @device_id, @event_type, @timestamp, @data, @created_at, @updated_at)`

	insertInboundQuery = `
       INSERT INTO "whatsapp-web".inbound 
	       (event_id, device_id, timestamp, data, created_at, updated_at) 
       VALUES (@event_id, @device_id, @timestamp, @data, @created_at, @updated_at)`
)

// SaveInboundEvent saves an inbound event to the inbound table
func (r *repo) SaveInbound(ctx context.Context, event *model.QueueEvent) error {
	data, err := util.MarshalToJSON(event.Data)
	if err != nil {
		return err
	}

	params := pgx.NamedArgs{
		"event_id":   event.EventID,
		"device_id":  event.DeviceID,
		"timestamp":  event.Timestamp,
		"data":       data,
		"created_at": util.NowUTC(),
		"updated_at": util.NowUTC(),
	}

	_, err = r.conn.Exec(ctx, insertInboundQuery, params)
	return err
}

func (r *repo) SaveEvent(ctx context.Context, event *model.QueueEvent) error {
	data, err := util.MarshalToJSON(event.Data)
	if err != nil {
		return err
	}

	params := pgx.NamedArgs{
		"event_id":   event.EventID,
		"device_id":  event.DeviceID,
		"event_type": event.EventType,
		"timestamp":  event.Timestamp,
		"data":       data,
		"created_at": util.NowUTC(),
		"updated_at": util.NowUTC(),
	}

	_, err = r.conn.Exec(ctx, insertEventQuery, params)
	return err
}
