package repository

import (
	"context"
	"eventhandler/model/entity"
	"eventhandler/util"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type WhatsAppAccountRepository interface {
	GetAccountBySenderJID(ctx context.Context, senderJID string) (*entity.WhatsAppAccount, error)
	UpdatePairingSuccess(ctx context.Context, req entity.WhatsAppAccountReq) error
	UpdateConnected(ctx context.Context, req entity.WhatsAppAccountReq) error
	UpdateDisconnected(ctx context.Context, req entity.WhatsAppAccountReq) error
}

type whatsappAccountRepo struct {
	pool *pgxpool.Pool
}

func NewWhatsAppAccountRepository(pool *pgxpool.Pool) WhatsAppAccountRepository {
	return &whatsappAccountRepo{
		pool: pool,
	}
}

const (
	getBySenderJIDQuery = `
		SELECT account_id, user_id, account_name, account_alias, phone_number, sender_jid,
			   session_data, connect_status, is_active, initiated_at, connected_at, 
			   disconnected_at, created_at, created_by
		FROM public.whatsapp_accounts 
		WHERE (sender_jid = @sender_jid OR account_id::text = @sender_jid) AND deleted_at IS NULL
		LIMIT 1
	`

	updatePairingSuccessQuery = `
		UPDATE public.whatsapp_accounts 
		SET phone_number = @phone_number, 
			sender_jid = @sender_jid, 
			connect_status = @connect_status, 
			connected_at = @connected_at,
			updated_at = @updated_at
		WHERE account_id = @account_id AND deleted_at IS NULL
	`

	updateConnectedQuery = `
		UPDATE public.whatsapp_accounts 
		SET connect_status = @connect_status, 
			connected_at = @connected_at,
			updated_at = @updated_at
		WHERE account_id = @account_id AND deleted_at IS NULL
	`

	updateDisconnectedQuery = `
		UPDATE public.whatsapp_accounts 
		SET connect_status = @connect_status, 
			disconnected_at = @disconnected_at,
			updated_at = @updated_at
		WHERE account_id = @account_id AND deleted_at IS NULL
	`
)

func (r *whatsappAccountRepo) GetAccountBySenderJID(ctx context.Context, senderJID string) (*entity.WhatsAppAccount, error) {
	params := pgx.NamedArgs{
		"sender_jid": senderJID,
	}

	var account entity.WhatsAppAccount
	err := r.pool.QueryRow(ctx, getBySenderJIDQuery, params).Scan(
		&account.AccountID,
		&account.UserID,
		&account.AccountName,
		&account.AccountAlias,
		&account.PhoneNumber,
		&account.SenderJID,
		&account.SessionData,
		&account.ConnectStatus,
		&account.IsActive,
		&account.InitiatedAt,
		&account.ConnectedAt,
		&account.DisconnectedAt,
		&account.CreatedAt,
		&account.CreatedBy,
	)
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (r *whatsappAccountRepo) UpdatePairingSuccess(ctx context.Context, req entity.WhatsAppAccountReq) error {
	params := pgx.NamedArgs{
		"account_id":     req.AccountID,
		"phone_number":   req.PhoneNumber,
		"sender_jid":     req.SenderJID,
		"connect_status": req.ConnectStatus,
		"connected_at":   req.ConnectedAt,
		"updated_at":     util.NowUTC(),
	}

	_, err := r.pool.Exec(ctx, updatePairingSuccessQuery, params)
	return err
}

func (r *whatsappAccountRepo) UpdateConnected(ctx context.Context, req entity.WhatsAppAccountReq) error {
	params := pgx.NamedArgs{
		"account_id":     req.AccountID,
		"connect_status": req.ConnectStatus,
		"connected_at":   req.ConnectedAt,
		"updated_at":     util.NowUTC(),
	}

	_, err := r.pool.Exec(ctx, updateConnectedQuery, params)
	return err
}

func (r *whatsappAccountRepo) UpdateDisconnected(ctx context.Context, req entity.WhatsAppAccountReq) error {
	params := pgx.NamedArgs{
		"account_id":      req.AccountID,
		"connect_status":  req.ConnectStatus,
		"disconnected_at": req.DisconnectedAt,
		"updated_at":      util.NowUTC(),
	}

	_, err := r.pool.Exec(ctx, updateDisconnectedQuery, params)
	return err
}
