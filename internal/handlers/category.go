package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/spalqui/habitattrack-api/internal/models"
	"github.com/spalqui/habitattrack-api/internal/services"
	"github.com/spalqui/habitattrack-api/pkg/utils"
)

type CategoryHandler struct {
	categoryService services.CategoryService
}

func NewCategoryHandler(categoryService services.CategoryService) *CategoryHandler {
	return &CategoryHandler{
		categoryService: categoryService,
	}
}

func (h *CategoryHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var category models.Category
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.categoryService.CreateCategory(r.Context(), &category); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.WriteJSONResponse(w, http.StatusCreated, category)
}

func (h *CategoryHandler) GetCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	category, err := h.categoryService.GetCategory(r.Context(), id)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusNotFound, err.Error())
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, category)
}

func (h *CategoryHandler) GetAllCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.categoryService.GetAllCategories(r.Context())
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, categories)
}

func (h *CategoryHandler) GetCategoriesByType(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	transactionType := models.TransactionType(vars["type"])

	categories, err := h.categoryService.GetCategoriesByType(r.Context(), transactionType)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, categories)
}

func (h *CategoryHandler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var category models.Category
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	category.ID = id
	if err := h.categoryService.UpdateCategory(r.Context(), &category); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, category)
}

func (h *CategoryHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := h.categoryService.DeleteCategory(r.Context(), id); err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
