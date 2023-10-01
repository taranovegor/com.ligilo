package service

import (
	"errors"
	kontrakto "github.com/taranovegor/com.kontrakto"
	"github.com/taranovegor/com.ligilo/internal/config"
	"github.com/taranovegor/com.ligilo/internal/domain"
	"gorm.io/gorm"
	"regexp"
)

type Validator struct {
	repository domain.LinkRepository
}

func NewValidator(
	repository domain.LinkRepository,
) *Validator {
	return &Validator{
		repository: repository,
	}
}

func (v Validator) ValidateToken(token string) kontrakto.ValidationResult {
	valid, _ := regexp.MatchString(config.LinkTokenRegex, token)
	result := kontrakto.ValidationResult{Success: valid}
	if !valid {
		result.Message = "Token is not valid"
	}

	if _, err := v.repository.GetByToken(token); !errors.Is(err, gorm.ErrRecordNotFound) {
		result.Success = false
		result.Message = "Token already in use"
	}

	return result
}

func (v Validator) ValidateLocation(location string) kontrakto.ValidationResult {
	valid, _ := regexp.MatchString(kontrakto.RegexHttpUrl, location)
	result := kontrakto.ValidationResult{Success: valid}
	if !valid {
		result.Message = "Link is not valid"
	}

	return result
}
