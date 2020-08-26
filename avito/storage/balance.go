package storage

import (
	"avito/db"
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
)

type BalanceStorageAPI interface {
	BalanceIncrease(tx pgx.Tx, userID uuid.UUID, sum int64) error
	BalanceDecrease(tx pgx.Tx, userID uuid.UUID, sum int64) error
	GetBalance(userID uuid.UUID) (int64, error)
	CountUsers(userID uuid.UUID) (int, error)
}

type balanceStorage struct {
	db *db.ConnDB
	ctx context.Context
}


func NewBalanceStorageAPI(connDB *db.ConnDB, ctx context.Context) BalanceStorageAPI {
	return &balanceStorage{
		db: connDB,
		ctx: ctx,
	}
}

func (c *balanceStorage) BalanceIncrease(tx pgx.Tx, userID uuid.UUID, sum int64) error {
	_, err := tx.Exec(c.ctx,"insert into balance (user_id, amount) values ($1, $2) on conflict (user_id) do update set amount = (select amount from balance where user_id = $1) + $2;", userID, sum)
	if err != nil {
		return err
	}

	return nil
}

func (c *balanceStorage) BalanceDecrease(tx pgx.Tx, userID uuid.UUID, sum int64) error {
	_, err := tx.Exec(c.ctx,"update balance set amount = (select amount from balance where user_id = $1) - $2 where user_id=$1;", userID, sum)
	if err != nil {
		return err
	}

	return nil
}

func (c *balanceStorage) GetBalance(userID uuid.UUID) (int64, error) {
	var result int64
	err := c.db.DB.QueryRow(c.ctx, "select amount from balance where user_id=$1", userID).Scan(&result)
	if err != nil {
		return 0, err
	}

	return result, nil
}

func (c *balanceStorage) CountUsers(userID uuid.UUID) (int, error) {
	var result int
	err := c.db.DB.QueryRow(c.ctx, "select count(user_id) from balance where user_id=$1", userID).Scan(&result)
	if err != nil {
		return 0, err
	}

	return result, nil
}