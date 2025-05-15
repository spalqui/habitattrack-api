package handlers

import "github.com/spalqui/habitattrack-api/types"

type CreatePropertyRequest struct {
	Number     string `json:"number" binding:"required"`
	StreetName string `json:"streetName" binding:"required"`
	Town       string `json:"town" binding:"required"`
	Postcode   string `json:"postcode" binding:"required"`
}

type UpdatePropertyRequest struct {
	Number     *string `json:"number,omitempty"`
	StreetName *string `json:"streetName,omitempty"`
	Town       *string `json:"town,omitempty"`
	Postcode   *string `json:"postcode,omitempty"`
}

func (r *CreatePropertyRequest) ToProperty() *types.Property {
	return &types.Property{
		Number:     r.Number,
		StreetName: r.StreetName,
		Town:       r.Town,
		Postcode:   r.Postcode,
	}
}

func (r *UpdatePropertyRequest) ApplyUpdates(p *types.Property) {
	if r.Number != nil {
		p.Number = *r.Number
	}
	if r.StreetName != nil {
		p.StreetName = *r.StreetName
	}
	if r.Town != nil {
		p.Town = *r.Town
	}
	if r.Postcode != nil {
		p.Postcode = *r.Postcode
	}
}
