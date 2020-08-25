package storage

import (
	"../db"
"context"
"github.com/jackc/pgx"
)

type StorageAPI interface {
	GetBalanceStorage() BalanceStorageAPI
	GetTransactionStorage() TransactionStorageAPI
	GetTransaction(ctx context.Context) (pgx.Tx, error)
}

type storageAPI struct {
	balanceStorage BalanceStorageAPI
	transactionStorage TransactionStorageAPI
	connDB db.ConnDB
}

func (s *storageAPI) GetTransaction(ctx context.Context) (pgx.Tx, error) {
	tx, err := s.connDB.DB.Begin(ctx)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func (s *storageAPI) GetBalanceStorage() BalanceStorageAPI {
	return s.balanceStorage
}

func (s *storageAPI) GetTransactionStorage() TransactionStorageAPI {
	return s.transactionStorage
}

func NewStorageAPI(connDB db.ConnDB, ctx context.Context) StorageAPI {
	return &storageAPI{
		balanceStorage: NewBalanceStorageAPI(connDB, ctx),
		transactionStorage: NewTransactionStorageAPI(connDB, ctx),
		connDB: connDB,
	}
}
