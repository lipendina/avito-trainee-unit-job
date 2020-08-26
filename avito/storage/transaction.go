package storage

import (
	"avito/db"
	"avito/dto"
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
)

type TransactionStorageAPI interface {
	GetTransactions(userID uuid.UUID, limit int, offset int) ([]dto.Transaction, error)
	WriteTransaction(tx pgx.Tx, userID uuid.UUID, sum int64) error
}

type transactionStorage struct {
	db *db.ConnDB
	ctx context.Context
}

func NewTransactionStorageAPI(connDB *db.ConnDB, ctx context.Context) TransactionStorageAPI {
	return &transactionStorage {
		db: connDB,
		ctx: ctx,
	}
}

func (t *transactionStorage) GetTransactions(userID uuid.UUID, limit int, offset int) ([]dto.Transaction, error) {

	rows, err := t.db.DB.Query(t.ctx, "select * from \"transaction\" where user_id=$1 order by created_at desc, change_balance asc limit $2 offset $3;", userID, limit, offset)
	if err != nil {
		return nil, err
	}

	result := make([]dto.Transaction, 0)
	for rows.Next() {
		var transaction dto.Transaction
		var money int64
		err := rows.Scan(&transaction.Id, &transaction.UserID, &money, &transaction.CreatedAt)
		if err != nil {
			return nil, err
		}

		transaction.ChangeBalance = &dto.Money{IntPart: money / 100, FracPart: money % 100}

		result = append(result, transaction)
	}

	return result, nil
}

func (t *transactionStorage) WriteTransaction(tx pgx.Tx, userID uuid.UUID, sum int64) error {
	_, err := tx.Exec(t.ctx,"insert into \"transaction\" (user_id, change_balance) values ($1, $2);", userID, sum)
	if err != nil {
		return err
	}

	return nil
}
