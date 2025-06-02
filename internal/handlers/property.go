package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/spalqui/habitattrack-api/internal/models"
	"github.com/spalqui/habitattrack-api/internal/services"
	"github.com/spalqui/habitattrack-api/pkg/utils"
)

type PropertyHandler struct {
	propertyService services.PropertyService
}

func NewPropertyHandler(propertyService services.PropertyService) *PropertyHandler {
	return &PropertyHandler{
		propertyService: propertyService,
	}
}

func (h *PropertyHandler) CreateProperty(w http.ResponseWriter, r *http.Request) {
	var property models.Property
	if err := json.NewDecoder(r.Body).Decode(&property); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.propertyService.CreateProperty(r.Context(), &property); err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.WriteJSONResponse(w, http.StatusCreated, property)
}

func (h *PropertyHandler) GetProperty(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	property, err := h.propertyService.GetProperty(r.Context(), id)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusNotFound, err.Error())
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, property)
}

func (h *PropertyHandler) GetAllProperties(w http.ResponseWriter, r *http.Request) {
	properties, err := h.propertyService.GetAllProperties(r.Context())
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, properties)
}

func (h *PropertyHandler) UpdateProperty(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var property models.Property
	if err := json.NewDecoder(r.Body).Decode(&property); err != nil {
		utils.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	property.ID = id
	if err := h.propertyService.UpdateProperty(r.Context(), &property); err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, property)
}

func (h *PropertyHandler) DeleteProperty(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := h.propertyService.DeleteProperty(r.Context(), id); err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
