package service

import (
	"github.com/matisiekpl/pansa-plan/internal/model"
	"github.com/matisiekpl/pansa-plan/internal/repository"
)

type ReportService interface {
	GenerateReport(plan model.Plan) ([]byte, error)
}

type reportService struct {
	reportRepository  repository.ReportRepository
	weatherRepository repository.WeatherRepository
}

func newReportService(reportRepository repository.ReportRepository, weatherRepository repository.WeatherRepository) ReportService {
	return reportService{reportRepository, weatherRepository}
}

func (r reportService) GenerateReport(plan model.Plan) ([]byte, error) {
	var weather []model.Weather
	for _, waypoint := range plan.Waypoints {
		if len(waypoint.Icao) == 4 {
			metar := r.weatherRepository.GetMETAR(waypoint.Icao)
			taf := r.weatherRepository.GetTAF(waypoint.Icao)
			if metar != "" {
				weather = append(weather, model.Weather{Name: waypoint.Icao, METAR: metar, TAF: taf})
			}
		}
	}

	return r.reportRepository.GenerateReport(plan, weather)
}
