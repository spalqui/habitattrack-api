package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/spalqui/habitattrack-api/internal/models"
	"github.com/spalqui/habitattrack-api/internal/services"
	"github.com/spalqui/habitattrack-api/pkg/utils"
)

type TransactionHandler struct {
	transactionService services.TransactionService
}

func NewTransactionHandler(transactionService services.TransactionService) *TransactionHandler {
	return &TransactionHandler{
		transactionService: transactionService,
	}
}

func (h *TransactionHandler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	var transaction models.Transaction
	if err := json.NewDecoder(r.Body).Decode(&transaction); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.transactionService.CreateTransaction(r.Context(), &transaction); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.WriteJSONResponse(w, http.StatusCreated, transaction)
}

func (h *TransactionHandler) GetTransaction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	transaction, err := h.transactionService.GetTransaction(r.Context(), id)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusNotFound, err.Error())
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, transaction)
}

func (h *TransactionHandler) GetAllTransactions(w http.ResponseWriter, r *http.Request) {
	transactions, err := h.transactionService.GetAllTransactions(r.Context())
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, transactions)
}

func (h *TransactionHandler) GetTransactionsByProperty(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	propertyID := vars["propertyId"]

	transactions, err := h.transactionService.GetTransactionsByProperty(r.Context(), propertyID)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, transactions)
}

func (h *TransactionHandler) UpdateTransaction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var transaction models.Transaction
	if err := json.NewDecoder(r.Body).Decode(&transaction); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	transaction.ID = id
	if err := h.transactionService.UpdateTransaction(r.Context(), &transaction); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, transaction)
}

func (h *TransactionHandler) DeleteTransaction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := h.transactionService.DeleteTransaction(r.Context(), id); err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
