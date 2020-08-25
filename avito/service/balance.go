package service

import (
	"../dto"
	"../storage"
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"golang.org/x/xerrors"
	"log"
	"net/http"
	"os"
)

// последний параметр в функциях - isInternal, для определения типа ошибки в handlers
type BalanceServiceAPI interface {
	CreditFundsRequest(creditFundsRequest dto.OperationRequest) (error, bool)
	WithdrawFundsRequest(withdrawFundsRequest dto.OperationRequest) (error, bool)
	TransferFundsRequest(transferFundsRequest dto.TransferFundsRequest) (error, bool)
	GetBalanceRequest(userID uuid.UUID, currency string) (*dto.Money, error, bool)
}

type balanceService struct {
	storage storage.StorageAPI
	ctx context.Context
	log *log.Logger
}

func NewBalanceServiceAPI(api storage.StorageAPI) BalanceServiceAPI {
	return &balanceService{
		storage: api,
		ctx: context.Background(),
		log: log.New(os.Stdout, "BALANCE-SERVICE: ", log.LstdFlags),
	}
}

func (b *balanceService) CreditFundsRequest(creditFundsRequest dto.OperationRequest) (error, bool) {
	b.log.Printf("Trying to increase balance of user %v", creditFundsRequest.UserId)

	if creditFundsRequest.Sum.FracPart  < 0 || creditFundsRequest.Sum.FracPart > 99 {
		return xerrors.Errorf("frac_part must be between 0 and 99"), false
	}

	sum := creditFundsRequest.Sum.IntPart * 100 + creditFundsRequest.Sum.FracPart
	if sum == 0 {
		return xerrors.Errorf("Sum cannot be 0"), false
	}

	tx, err := b.storage.GetTransaction(b.ctx)
	if err != nil {
		b.log.Printf("Error while create transaction, reason: %+v", err)
		return xerrors.Errorf("System error. Contact support"), true
	}

	err = b.storage.GetBalanceStorage().BalanceIncrease(tx, creditFundsRequest.UserId, sum)
	if err != nil {
		b.log.Printf("Error while increase balance in DB, reason: %v", err)
		tx.Rollback(b.ctx)
		return xerrors.Errorf("System error. Contact support"), true
	}

	err = b.storage.GetTransactionStorage().WriteTransaction(tx, creditFundsRequest.UserId, sum)
	if err != nil {
		b.log.Printf("Error while write transaction in DB, reason: %v", err)
		tx.Rollback(b.ctx)
		return xerrors.Errorf("System error. Contact support"), true
	}

	err = tx.Commit(b.ctx)
	if err != nil {
		b.log.Printf("Error while commit transaction, reason: %+v", err)
		return xerrors.Errorf("System error. Contact support"), true
	}

	return nil, false
}

func (b *balanceService) WithdrawFundsRequest(withdrawFundsRequest dto.OperationRequest) (error, bool) {
	b.log.Printf("Trying to decrease balance of user %v", withdrawFundsRequest.UserId)

	if withdrawFundsRequest.Sum.FracPart  < 0 || withdrawFundsRequest.Sum.FracPart > 99 {
		return xerrors.Errorf("frac_part must be between 0 and 99"), false
	}

	sum := withdrawFundsRequest.Sum.IntPart * 100 + withdrawFundsRequest.Sum.FracPart
	if sum == 0 {
		return xerrors.Errorf("Sum cannot be 0"), false
	}

	count, err := b.storage.GetBalanceStorage().CountUsers(withdrawFundsRequest.UserId)
	if err != nil {
		b.log.Printf("Error while count users in DB, reason: %v", err)
		return xerrors.Errorf("System error. Contact support"), true
	}

	if count != 1 {
		return xerrors.Errorf("User does not exist"), false
	}

	balance, err := b.storage.GetBalanceStorage().GetBalance(withdrawFundsRequest.UserId)
	if err != nil {
		b.log.Printf("Error while get balance from DB, reason: %v", withdrawFundsRequest)
		return xerrors.Errorf("System error. Contact support"), true
	}

	if balance < sum {
		return xerrors.Errorf("You have not enough funds to complete this operation"), false
	}

	tx, err := b.storage.GetTransaction(b.ctx)
	if err != nil {
		b.log.Printf("Error while create transaction, reason: %+v", err)
		return xerrors.Errorf("System error. Contact support"), true
	}

	err = b.storage.GetBalanceStorage().BalanceDecrease(tx, withdrawFundsRequest.UserId, sum)
	if err != nil {
		b.log.Printf("Error while decrease balance in DB, reason: %v", err)
		tx.Rollback(b.ctx)
		return xerrors.Errorf("System error. Contact support"), true
	}

	err = b.storage.GetTransactionStorage().WriteTransaction(tx, withdrawFundsRequest.UserId, -sum)
	if err != nil {
		b.log.Printf("Error while write transaction in DB, reason: %v", err)
		tx.Rollback(b.ctx)
		return xerrors.Errorf("System error. Contact support"), true
	}

	err = tx.Commit(b.ctx)
	if err != nil {
		b.log.Printf("Error while commit transaction, reason: %+v", err)
		return xerrors.Errorf("System error. Contact support"), true
	}

	return nil, false
}

func (b *balanceService) TransferFundsRequest(transferFundsRequest dto.TransferFundsRequest) (error, bool) {
	b.log.Printf("Trying to transfer funds from user %v to user %v", transferFundsRequest.IdSender, transferFundsRequest.IdReceiver)

	if transferFundsRequest.Sum.FracPart  < 0 || transferFundsRequest.Sum.FracPart > 99 {
		return xerrors.Errorf("frac_part must be between 0 and 99"), false
	}

	sum := transferFundsRequest.Sum.IntPart * 100 + transferFundsRequest.Sum.FracPart
	if transferFundsRequest.IdReceiver == transferFundsRequest.IdSender {
		return xerrors.Errorf("ReceiverID and senderID cannot be equal"), false
	}

	if sum == 0 {
		return xerrors.Errorf("Sum cannot be 0"), false
	}

	balance, err := b.storage.GetBalanceStorage().GetBalance(transferFundsRequest.IdSender)
	if err != nil {
		b.log.Printf("Error while get balance from DB, reason: %v", err)
		return xerrors.Errorf("System error. Contact support"), true
	}

	if balance < sum {
		return xerrors.Errorf("You have not enough funds to complete this operation"), false
	}

	tx, err := b.storage.GetTransaction(b.ctx)
	if err != nil {
		b.log.Printf("Error while create transaction, reason: %+v", err)
		return xerrors.Errorf("System error. Contact support"), true
	}

	err = b.storage.GetBalanceStorage().BalanceDecrease(tx, transferFundsRequest.IdSender, sum)
	if err != nil {
		b.log.Printf("Error while decrease balance in DB, reason: %v", err)
		tx.Rollback(b.ctx)
		return xerrors.Errorf("System error. Contact support"), true
	}

	err = b.storage.GetTransactionStorage().WriteTransaction(tx, transferFundsRequest.IdSender, -sum)
	if err != nil {
		b.log.Printf("Error while write transaction in DB, reason: %v", err)
		tx.Rollback(b.ctx)
		return xerrors.Errorf("System error. Contact support"), true
	}

	err = b.storage.GetBalanceStorage().BalanceIncrease(tx, transferFundsRequest.IdReceiver, sum)
	if err != nil {
		b.log.Printf("Error while increase balance in DB, reason: %v", err)
		tx.Rollback(b.ctx)
		return xerrors.Errorf("System error. Contact support"), true
	}

	err = b.storage.GetTransactionStorage().WriteTransaction(tx, transferFundsRequest.IdReceiver, sum)
	if err != nil {
		b.log.Printf("Error while write transaction in DB, reason: %v", err)
		tx.Rollback(b.ctx)
		return xerrors.Errorf("System error. Contact support"), true
	}

	err = tx.Commit(b.ctx)
	if err != nil {
		b.log.Printf("Error while commit transaction, reason: %+v", err)
		return xerrors.Errorf("System error. Contact support"), true
	}

	return nil, false
}

func (b *balanceService) GetBalanceRequest(userID uuid.UUID, currency string) (*dto.Money, error, bool) {
	b.log.Printf("Trying to get balance of user %v", userID)

	count, err := b.storage.GetBalanceStorage().CountUsers(userID)
	if err != nil {
		b.log.Printf("Error while count users in DB, reason: %v", err)
		return nil, xerrors.Errorf("System error. Contact support"), true
	}

	if count != 1 {
		return nil, xerrors.Errorf("User does not exist"), false
	}

	balance, err := b.storage.GetBalanceStorage().GetBalance(userID)
	if err != nil {
		b.log.Printf("Error while get balance from DB, reason: %v", err)
		return nil, xerrors.Errorf("System error. Contact support"), true
	}


	if currency != "" {
		cur, err, isUserError := GetCurrencyRequest(currency)
		if err != nil {
			b.log.Printf("Error while get currency, reason: %v", err)
			if isUserError {
				return nil, err, false
			}
			return nil, xerrors.Errorf("System error. Contact support"), true
		}
		s := float64(balance) * cur
		return &dto.Money{IntPart: int64(s) / 100, FracPart: int64(s) % 100}, nil, false

	}

	return &dto.Money{IntPart: balance / 100, FracPart: balance % 100}, nil, false
}

func GetCurrencyRequest(value string) (float64, error, bool) {
	r, err := http.Get("https://api.exchangeratesapi.io/latest?base=RUB")
	if err != nil {
		return 0, err, false
	}
	defer r.Body.Close()

	var result dto.CurrencyRates
	err = json.NewDecoder(r.Body).Decode(&result)
	if err != nil {
		return 0, err, false
	}

	if _, ok := result.Rates[value]; ok == false {
		return 0, xerrors.Errorf("Currency is not exist"), true
	}

	return result.Rates[value], nil, false
}
