package service

import "github.com/matisiekpl/pansa-plan/internal/repository"

type Services interface {
	Report() ReportService
	Publication() PublicationService
	Notam() NotamService
}

type services struct {
	report      ReportService
	publication PublicationService
	notam       NotamService
}

func NewServices(repositories repository.Repositories) Services {
	publicationService := newPublicationService(repositories.Publication())
	notamService := newNotamService(repositories.Notam())
	reportService := newReportService(repositories.Report())
	return &services{
		publication: publicationService,
		report:      reportService,
		notam:       notamService,
	}
}

func (s services) Report() ReportService {
	return s.report
}

func (s services) Publication() PublicationService {
	return s.publication
}

func (s services) Notam() NotamService {
	return s.notam
}
