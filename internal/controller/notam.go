package controller

import (
	"github.com/labstack/echo/v4"
	"github.com/matisiekpl/pansa-plan/internal/service"
	"net/http"
)

type NotamController interface {
	Index(c echo.Context) error
}

type notamController struct {
	notamService service.NotamService
}

func newNotamController(notamService service.NotamService) NotamController {
	return &notamController{notamService}
}

func (n notamController) Index(c echo.Context) error {
	notam, err := n.notamService.Index(c.Param("icao"))
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, notam)
}
