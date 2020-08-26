package main

import (
	"avito/handlers"
	"avito/storage"
	"avito/service"
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"avito/config"
	"avito/db"
)

func main() {

	applicationConfig, err := config.ParseConfig()
	if err != nil {
		log.Fatalf("Cannot parse config: %+v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	pgConn, err := db.NewConnectToPG(&applicationConfig.DB, ctx)
	if err != nil {
		log.Fatalf("Cannot connect to DB, reason: %v", err)
	}

	storageAPI := storage.NewStorageAPI(pgConn, ctx)
	serviceAPI := service.NewServiceAPI(storageAPI)

	a := handlers.NewHandlers(serviceAPI)

	r := mux.NewRouter()
	// зачисление денежных средств
	r.HandleFunc("/balance/credit", a.CreditFundsHandler)
	// списание денежных средств
	r.HandleFunc("/balance/withdraw", a.WithdrawFundsHandler)
	// перевод денежных средств другому пользователю
	r.HandleFunc("/balance/transfer", a.TransferFundsHandler)
	// получение текущего баланса
	r.HandleFunc("/balance/get", a.GetBalanceHandler)
	// получение
	r.HandleFunc("/balance/transactions", a.GetTransactionsHandler)
	http.Handle("/", r)

	fmt.Println("Server is listening...")
	http.ListenAndServe(fmt.Sprintf(":%d", applicationConfig.HTTPPort), nil)
}
