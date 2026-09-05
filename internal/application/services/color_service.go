package services

import (
	"errors"
	"regexp"
	"strings"

	"github.com/mahditd/zarrine-baft-backend/internal/domain/models"
	"github.com/mahditd/zarrine-baft-backend/internal/domain/repositories"
)

type ColorService struct {
	colorRepository repositories.ColorRepository
}

func NewColorService(
	colorRepository repositories.ColorRepository,
) *ColorService {

	return &ColorService{
		colorRepository: colorRepository,
	}
}

type CreateColorInput struct {
	NameFA  string
	NameEN  string
	HexCode string
}

func (s *ColorService) Create(
	input CreateColorInput,
) (*models.Color, error) {

	input.NameFA = normalizeColorName(input.NameFA)
	input.NameEN = normalizeColorName(input.NameEN)
	input.HexCode = strings.TrimSpace(input.HexCode)

	if input.NameFA == "" {
		return nil, errors.New("persian name is required")
	}

	if input.NameEN == "" {
		return nil, errors.New("english name is required")
	}

	if !isValidHexCode(input.HexCode) {
		return nil, errors.New("invalid hex code")
	}

	existingFA, err := s.colorRepository.FindByNameFA(
		input.NameFA,
	)

	if err == nil && existingFA != nil {
		return nil, errors.New("persian color name already exists")
	}

	existingEN, err := s.colorRepository.FindByNameEN(
		input.NameEN,
	)

	if err == nil && existingEN != nil {
		return nil, errors.New("english color name already exists")
	}

	color := &models.Color{
		NameFA:  input.NameFA,
		NameEN:  input.NameEN,
		HexCode: input.HexCode,
	}

	err = s.colorRepository.Create(color)

	if err != nil {
		return nil, err
	}

	return color, nil
}

func (s *ColorService) GetAll() ([]models.Color, error) {

	return s.colorRepository.FindAll()

}

func normalizeColorName(name string) string {

	name = strings.TrimSpace(name)

	name = strings.Join(
		strings.Fields(name),
		" ",
	)

	return name
}

func isValidHexCode(hex string) bool {

	pattern := `^#[0-9A-Fa-f]{6}$`

	matched, _ := regexp.MatchString(
		pattern,
		hex,
	)

	return matched
}
