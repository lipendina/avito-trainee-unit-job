package dto

import (
	"fmt"
	"github.com/google/uuid"
)

type OperationRequest struct {
	UserId uuid.UUID `json:"user_id"`
	Sum *Money `json:"amount"`
}

type TransferFundsRequest struct {
	IdSender uuid.UUID `json:"sender_id"`
	IdReceiver uuid.UUID `json:"receiver_id"`
	Sum *Money `json:"amount"`
}

type GetBalanceResponse struct {
	Sum *Money `json:"amount"`
}

type GetTransactionsResponse struct {
	Transactions []Transaction
}

type Transaction struct {
	Id uuid.UUID `json:"id"`
	UserID uuid.UUID `json:"user_id"`
	ChangeBalance *Money `json:"change_balance"`
	CreatedAt string `json:"created_at"`
}

type ErrorResponse struct {
	Error string
}

type Money struct {
	IntPart int64 `json:"int_part"`
	FracPart int64 `json:"frac_part"`
}

type CurrencyRates struct {
	Rates map[string]float64 `json:"rates"`
	Base string `json:"base"`
	Date string `json:"date"`
}

func (r OperationRequest) String() string {
	return fmt.Sprintf("{User ID: %v, sum: %v}", r.UserId, r.Sum)
}

func (r TransferFundsRequest) String() string {
	return fmt.Sprintf("{Sender ID: %v, receiver ID: %v, sum: %v}", r.IdSender, r.IdReceiver, r.Sum)
}

func (r ErrorResponse) String() string {
	return fmt.Sprintf("Error: %s", r.Error)
}

func (r GetBalanceResponse) String() string {
	return fmt.Sprintf("Balance: %d.%d", r.Sum.IntPart, r.Sum.FracPart)
}

func (r GetTransactionsResponse) String() string {
	return fmt.Sprintf("Transactions: %v", r.Transactions)
}

func (r Transaction) String() string {
	return fmt.Sprintf("{ID: %v, user id: %v, change: %v, created at: %v}", r.Id, r.UserID, r.ChangeBalance, r.CreatedAt)
}

func (r Money) String() string {
	return fmt.Sprintf("%d.%d", r.IntPart, r.FracPart)
}


