package dto

import (
	"github.com/mahditd/zarrine-baft-backend/internal/domain/models"
)

type ColorResponse struct {
	ID      uint   `json:"id"`
	NameFA  string `json:"name_fa"`
	NameEN  string `json:"name_en"`
	HexCode string `json:"hex_code"`
}

func FromColor(
	color *models.Color,
) ColorResponse {

	return ColorResponse{
		ID:      color.ID,
		NameFA:  color.NameFA,
		NameEN:  color.NameEN,
		HexCode: color.HexCode,
	}
}
