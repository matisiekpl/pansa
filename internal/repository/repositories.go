package repository

type Repositories interface {
	Publication() PublicationRepository
	Notam() NotamRepository
	Report() ReportRepository
	Weather() WeatherRepository
}

type repositories struct {
	publication PublicationRepository
	notam       NotamRepository
	report      ReportRepository
	weather     WeatherRepository
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
func (r repositories) Weather() WeatherRepository {
	return r.weather
}

func NewRepositories() Repositories {
	publicationRepository := newPublicationRepository()
	notamRepository := newNotamRepository()
	reportRepository := newReportRepository()
	weatherRepository := newWeatherRepository()
	return &repositories{
		publication: publicationRepository,
		notam:       notamRepository,
		report:      reportRepository,
		weather:     weatherRepository,
	}
}
