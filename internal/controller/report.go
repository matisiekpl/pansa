package controller

import (
	"github.com/labstack/echo/v4"
	"github.com/matisiekpl/pansa-plan/internal/model"
	"github.com/matisiekpl/pansa-plan/internal/service"
	"net/http"
)

type ReportController interface {
	Generate(c echo.Context) error
}

type reportController struct {
	reportService service.ReportService
}

func (r reportController) Generate(c echo.Context) error {
	var plan model.Plan
	err := c.Bind(&plan)
	if err != nil {
		return err
	}
	b, err := r.reportService.GenerateReport(plan)
	if err != nil {
		return err
	}
	return c.Blob(http.StatusOK, "application/pdf", b)
}

func newReportController(reportService service.ReportService) ReportController {
	return reportController{reportService}
}
