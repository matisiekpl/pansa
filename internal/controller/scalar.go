package controller

import (
	"embed"
	"github.com/labstack/echo/v4"
	"net/http"
)

//go:embed resources/openapi.yaml resources/scalar.html
var resources embed.FS

type ScalarController interface {
	Spec(c echo.Context) error
	UI(c echo.Context) error
}

type scalarController struct{}

func newScalarController() ScalarController {
	return &scalarController{}
}

func (s scalarController) Spec(c echo.Context) error {
	data, err := resources.ReadFile("resources/openapi.yaml")
	if err != nil {
		return err
	}
	return c.Blob(http.StatusOK, "application/yaml", data)
}

func (s scalarController) UI(c echo.Context) error {
	data, err := resources.ReadFile("resources/scalar.html")
	if err != nil {
		return err
	}
	return c.HTML(http.StatusOK, string(data))
}
