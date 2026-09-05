package dto

import (
	"github.com/mahditd/zarrine-baft-backend/internal/domain/models"
)

type MaterialResponse struct {
	ID     uint   `json:"id"`
	NameFA string `json:"name_fa"`
	NameEN string `json:"name_en"`
}

func FromMaterial(material *models.Material) MaterialResponse {

	return MaterialResponse{
		ID:     material.ID,
		NameFA: material.NameFA,
		NameEN: material.NameEN,
	}
}
