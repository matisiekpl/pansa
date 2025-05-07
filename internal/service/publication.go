package service

import (
	"github.com/matisiekpl/pansa-plan/internal/model"
	"github.com/matisiekpl/pansa-plan/internal/repository"
)

type PublicationService interface {
	Index(language model.Language) []model.Publication
}

type publicationService struct {
	publicationRepository repository.PublicationRepository
}

func newPublicationService(publicationRepository repository.PublicationRepository) PublicationService {
	return &publicationService{publicationRepository}
}

func (p publicationService) Index(language model.Language) []model.Publication {
	return p.publicationRepository.Index(language)
}
