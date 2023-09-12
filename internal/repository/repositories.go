package repository

type Repositories interface {
	Publication() PublicationRepository
	Notam() NotamRepository
	Report() ReportRepository
}

type repositories struct {
	publication PublicationRepository
	notam       NotamRepository
	report      ReportRepository
}

func (r repositories) Publication() PublicationRepository {
	return r.publication
}

func (r repositories) Notam() NotamRepository {
	return r.notam
}

func (r repositories) Report() ReportRepository {
	return r.report
}

func NewRepositories() Repositories {
	publicationRepository := newPublicationRepository()
	notamRepository := newNotamRepository()
	reportRepository := newReportRepository()
	return &repositories{
		publication: publicationRepository,
		notam:       notamRepository,
		report:      reportRepository,
	}
}
