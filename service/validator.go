package service

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"slices"
	"strings"
)

type Validator interface {
	GetValid(URIs []string) ([]string, error)
}

type ValidatorImpl struct{}

func NewValidator() Validator {
	return &ValidatorImpl{}
}

const URIPattern = "https?:\\/(?:\\/.+\\/)+(?:.+\\..+)"

func (v *ValidatorImpl) GetValid(URIs []string) ([]string, error) {
	var errs []string
	var validUris []string
	for _, URI := range URIs {
		matched, err := regexp.MatchString(URIPattern, URI)
		if err != nil || !matched {
			errs = append(errs, fmt.Sprintf("\"%s\" didn't pass validation and will be ignored", URI))
		} else if slices.Contains(validUris, URI) {
			errs = append(errs, fmt.Sprintf("\"%s\" already exists in queue and will be ignored", URI))
		} else {
			validUris = append(validUris, URI)
		}
	}

	if len(validUris) == 0 {
		log.Fatal("No valid URIs found")
	}

	if len(errs) != 0 {
		return validUris, errors.New(strings.Join(errs, "\n"))
	}

	return validUris, nil
}
