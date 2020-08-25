package handlers

import (
	"../dto"
	"../service"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"log"
	"net/http"
	"os"
	"strconv"
)

type Handlers interface {
	CreditFundsHandler(w http.ResponseWriter, r *http.Request)
	WithdrawFundsHandler(w http.ResponseWriter, r *http.Request)
	TransferFundsHandler(w http.ResponseWriter, r *http.Request)
	GetBalanceHandler(w http.ResponseWriter, r *http.Request)
	GetTransactionsHandler(w http.ResponseWriter, r *http.Request)
}

type handlers struct {
	service service.ServiceAPI
	log *log.Logger
}

func NewHandlers(api service.ServiceAPI) Handlers {
	return &handlers{
		service: api,
		log: log.New(os.Stdout, "CONTROLLER: ", log.LstdFlags),
	}
}

func (h *handlers) CreditFundsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var creditFundsRequest dto.OperationRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	err := dec.Decode(&creditFundsRequest)

	if err != nil {
		h.log.Printf("Error while parse creditFundsRequest, reason: %v", err)
		response := &dto.ErrorResponse{Error: "Cannot parse request body"}
		sendResponse(http.StatusBadRequest, response, w)
		return
	}
	h.log.Printf("Received creditFundsRequest: %v", creditFundsRequest)

	err, bool := h.service.GetBalanceService().CreditFundsRequest(creditFundsRequest)
	if err != nil {
		h.log.Printf("Error while do creditFundsRequest, reason: %v", err)
		response := &dto.ErrorResponse{Error: err.Error()}
		sendResponse(getErrorStatus(bool), response, w)
		return
	}

	response := fmt.Sprintf("OK")
	h.log.Printf("Funds have been successfully credited to the account of user with id %v", creditFundsRequest.UserId)
	sendResponse(http.StatusOK, response, w)
}

func (h *handlers) WithdrawFundsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var withdrawFundsRequest dto.OperationRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	err := dec.Decode(&withdrawFundsRequest)

	if err != nil {
		h.log.Printf("Error while parse withdrawFundsRequest, reason: %v", err)
		response := &dto.ErrorResponse{Error: "Cannot parse request body"}
		sendResponse(http.StatusBadRequest, response, w)
		return
	}
	h.log.Printf("Received withdrawFundsRequest: %v", withdrawFundsRequest)

	err, bool := h.service.GetBalanceService().WithdrawFundsRequest(withdrawFundsRequest)
	if err != nil {
		h.log.Printf("Error while do withdrawFundsRequest, reason: %v", err)
		response := &dto.ErrorResponse{Error: err.Error()}
		sendResponse(getErrorStatus(bool), response, w)
		return
	}

	response := fmt.Sprintf("OK")
	h.log.Printf("Funds have been successfully withdraw from the account of user %v", withdrawFundsRequest.UserId)
	sendResponse(http.StatusOK, response, w)
}

func (h *handlers) TransferFundsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var transferFundsRequest dto.TransferFundsRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	err := dec.Decode(&transferFundsRequest)

	if err != nil {
		h.log.Printf("Error while parse transferFundsRequest, reason: %v", err)
		response := &dto.ErrorResponse{Error: "Cannot parse request body"}
		sendResponse(http.StatusBadRequest, response, w)
		return
	}
	h.log.Printf("Received transferFundsRequest: %v", transferFundsRequest)

	err, bool := h.service.GetBalanceService().TransferFundsRequest(transferFundsRequest)
	if err != nil {
		h.log.Printf("Error while do transferFundsRequest, reason: %v", err)
		response := &dto.ErrorResponse{Error: err.Error()}
		sendResponse(getErrorStatus(bool), response, w)
		return
	}

	response := fmt.Sprintf("OK")
	h.log.Printf("Funds have been successfully transfer from user %v to user %v", transferFundsRequest.IdSender, transferFundsRequest.IdReceiver)
	sendResponse(http.StatusOK, response, w)

}

func (h *handlers) GetBalanceHandler(w http.ResponseWriter, r *http.Request) {

	uID := r.URL.Query().Get("user_id")
	if uID == "" {
		h.log.Printf("Error while parse value of user_id")
		response := &dto.ErrorResponse{Error: "Unknown user_id"}
		sendResponse(http.StatusBadRequest, response, w)
		return
	}

	userID, err := uuid.Parse(uID)
	if err != nil {
		h.log.Printf("Error while convert userID from string to uuid.UUID")
		response := &dto.ErrorResponse{Error: "System error. Contact support"}
		sendResponse(http.StatusBadRequest, response, w)
		return
	}

	currency := r.URL.Query().Get("currency")

	result, err, isInternal := h.service.GetBalanceService().GetBalanceRequest(userID, currency)
	if err != nil {
		h.log.Printf("Error while do creditFundsRequest, reason: %v", err)
		response := &dto.ErrorResponse{Error: err.Error()}
		sendResponse(getErrorStatus(isInternal), response, w)
		return
	}

	response := &dto.GetBalanceResponse{Sum: result}
	h.log.Printf("Send response: %v", response)
	sendResponse(http.StatusOK, fmt.Sprintf("Balance: %d.%d", response.Sum.IntPart, response.Sum.FracPart), w)
}

func (h *handlers) GetTransactionsHandler(w http.ResponseWriter, r *http.Request) {

	uID := r.URL.Query().Get("user_id")
	if uID == "" {
		h.log.Printf("Error while parse value of user_id")
		response := &dto.ErrorResponse{Error: "Unknown user_id"}
		sendResponse(http.StatusBadRequest, response, w)
		return
	}

	userID, err := uuid.Parse(uID)
	if err != nil {
		h.log.Printf("Error while convert userID from string to uuid.UUID")
		response := &dto.ErrorResponse{Error: "System error. Contact support"}
		sendResponse(http.StatusBadRequest, response, w)
	}

	l := r.URL.Query().Get("limit")
	var limit int
	if len(l) == 0 {
		limit = 100
	} else {
		limit, err = strconv.Atoi(l)
		if err != nil {
			h.log.Printf("Error while parse value of limit")
			response := &dto.ErrorResponse{Error: "Incorrect value of limit"}
			sendResponse(http.StatusBadRequest, response, w)
		}
	}

	ofs := r.URL.Query().Get("offset")
	var offset int
	if len(ofs) == 0 {
		offset = 0
	} else {
		offset, err = strconv.Atoi(ofs)
		if err != nil {
			h.log.Printf("Error while parse value of offset")
			response := &dto.ErrorResponse{Error: "Incorrect value of offset"}
			sendResponse(http.StatusBadRequest, response, w)
		}
	}

	rows, err, isInternal := h.service.GetTransactionService().GetTransactionsRequest(userID, limit, offset)
	if err != nil {
		h.log.Printf("Error while do getTransactionsRequest, reason: %v", err)
		response := &dto.ErrorResponse{Error: err.Error()}
		sendResponse(getErrorStatus(isInternal), response, w)
		return
	}

	response := &dto.GetTransactionsResponse{Transactions: rows}
	h.log.Printf("Send response: %v", response)
	sendResponse(http.StatusOK, response, w)
}
