package controller

import (
	"github.com/labstack/echo/v4"
	"github.com/matisiekpl/pansa-plan/internal/model"
	"github.com/matisiekpl/pansa-plan/internal/service"
	"net/http"
)

type PublicationController interface {
	Index(c echo.Context) error
}

type publicationController struct {
	publicationService service.PublicationService
}

func newPublicationController(publicationService service.PublicationService) PublicationController {
	return publicationController{publicationService}
}

func (p publicationController) Index(c echo.Context) error {
	language := c.QueryParam("language")
	if language == "" {
		language = string(model.LanguagePolish)
	}
	return c.JSON(http.StatusOK, p.publicationService.Index(model.Language(language)))
}
