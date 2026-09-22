package controller

import "github.com/matisiekpl/pansa-plan/internal/service"

type Controllers interface {
	Notam() NotamController
	Report() ReportController
	Publication() PublicationController
	Scalar() ScalarController
}

type controllers struct {
	notamController       NotamController
	reportController      ReportController
	publicationController PublicationController
	scalarController      ScalarController
}

func NewControllers(services service.Services) Controllers {
	notamController := newNotamController(services.Notam())
	reportController := newReportController(services.Report())
	publicationController := newPublicationController(services.Publication())
	scalarController := newScalarController()
	return &controllers{
		notamController:       notamController,
		reportController:      reportController,
		publicationController: publicationController,
		scalarController:      scalarController,
	}
}

func (c controllers) Notam() NotamController {
	return c.notamController
}

func (c controllers) Report() ReportController {
	return c.reportController
}

func (c controllers) Publication() PublicationController {
	return c.publicationController
}

func (c controllers) Scalar() ScalarController {
	return c.scalarController
}
