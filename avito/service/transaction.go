package service

import (
	"avito/dto"
	"avito/storage"
	"context"
	"github.com/google/uuid"
	"golang.org/x/xerrors"
	"log"
	"os"
)

// последний параметр в функциях - isInternal, для определения типа ошибки в handlers
type TransactionServiceAPI interface {
	GetTransactionsRequest(userID uuid.UUID, limit int, offset int) ([]dto.Transaction, error, bool)
}

type transactionService struct {
	storage storage.StorageAPI
	ctx context.Context
	log *log.Logger
}

func NewTransactionServiceAPI(api storage.StorageAPI) TransactionServiceAPI {
	return &transactionService{
		storage: api,
		ctx: context.Background(),
		log: log.New(os.Stdout, "TRANSACTION-SERVICE: ", log.LstdFlags),
	}
}

func (t *transactionService) GetTransactionsRequest(userID uuid.UUID, limit int, offset int) ([]dto.Transaction, error, bool) {
	t.log.Printf("Trying get transactions of user %v", userID)

	count, err := t.storage.GetBalanceStorage().CountUsers(userID)
	if err != nil {
		t.log.Printf("Error while count users in DB, reason: %v", err)
		return nil, xerrors.Errorf("System error. Contact support"), true
	}

	if count != 1 {
		return nil, xerrors.Errorf("User does not exist"), false
	}

	rows, err := t.storage.GetTransactionStorage().GetTransactions(userID, limit, offset)
	if err != nil {
		t.log.Printf("Error while get transactions from DB, reason: %v", err)
		return nil, xerrors.Errorf("System error. Contact support"), true
	}

	return rows, nil, false
}