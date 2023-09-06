package main

import (
	"github.com/labstack/echo/v4"
	"github.com/matisiekpl/notam"
	"net/http"
	"strings"
)

func fetchNotam(c echo.Context) error {
	items, err := notam.Fetch(strings.ToUpper(c.Param("icao")))
	if err != nil {
		return c.JSON(http.StatusOK, []string{})
	}
	return c.JSON(http.StatusOK, items)
}
