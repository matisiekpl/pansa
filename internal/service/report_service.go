package service

import (
	"github.com/matisiekpl/pansa-plan/internal/model"
	"github.com/matisiekpl/pansa-plan/internal/repository"
)

type ReportService interface {
	GenerateReport(plan model.Plan) ([]byte, error)
}

type reportService struct {
	reportRepository repository.ReportRepository
}

func newReportService(reportRepository repository.ReportRepository) ReportService {
	return reportService{reportRepository}
}

func (r reportService) GenerateReport(plan model.Plan) ([]byte, error) {
	return r.reportRepository.GenerateReport(plan)
}
