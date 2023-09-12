package service

import (
	"github.com/matisiekpl/pansa-plan/internal/model"
	"github.com/matisiekpl/pansa-plan/internal/repository"
)

type PublicationService interface {
	Index() []model.Publication
}

type publicationService struct {
	publicationRepository repository.PublicationRepository
}

func newPublicationService(publicationRepository repository.PublicationRepository) PublicationService {
	return &publicationService{publicationRepository}
}

func (p publicationService) Index() []model.Publication {
	return p.publicationRepository.Index()
}
