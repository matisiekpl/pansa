package service

import "github.com/matisiekpl/pansa-plan/internal/repository"

type NotamService interface {
	Index(icao string) ([]string, error)
}

type notamService struct {
	notamRepository repository.NotamRepository
}

func newNotamService(notamRepository repository.NotamRepository) *notamService {
	return &notamService{notamRepository: notamRepository}
}

func (s *notamService) Index(icao string) ([]string, error) {
	return s.notamRepository.Index(icao)
}
