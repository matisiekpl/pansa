package repository

import (
	"github.com/matisiekpl/notam"
	"strings"
)

type NotamRepository interface {
	Index(icao string) ([]string, error)
}

type notamRepository struct {
}

func newNotamRepository() NotamRepository {
	return &notamRepository{}
}

func (n notamRepository) Index(icao string) ([]string, error) {
	return notam.Fetch(strings.ToUpper(icao))
}
