package service

import (
	"../storage"
)

type ServiceAPI interface {
	GetBalanceService() BalanceServiceAPI
	GetTransactionService() TransactionServiceAPI
}

type serviceAPI struct {
	balanceServiceAPI BalanceServiceAPI
	transactionServiceAPI TransactionServiceAPI
}

func NewServiceAPI(api storage.StorageAPI) ServiceAPI {
	return &serviceAPI{
		balanceServiceAPI: NewBalanceServiceAPI(api),
		transactionServiceAPI: NewTransactionServiceAPI(api),
	}
}

func (s *serviceAPI) GetBalanceService() BalanceServiceAPI {
	return s.balanceServiceAPI
}

func (s *serviceAPI) GetTransactionService() TransactionServiceAPI {
	return s.transactionServiceAPI
}