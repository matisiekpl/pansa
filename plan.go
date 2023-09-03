package main

import (
	"github.com/desertbit/fillpdf"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"strings"
)

func generatePlan(c echo.Context) error {
	filename := "plan" + uuid.New().String() + ".pdf"
	logrus.Infof("creating flight plan for %s", strings.ToUpper(c.FormValue("aircraft identyfication")))
	form := fillpdf.Form{
		"aircraft identyfication":      strings.ToUpper(c.FormValue("aircraft identyfication")),
		"flight rules":                 strings.ToUpper(c.FormValue("flight rules")),
		"type of flight":               strings.ToUpper(c.FormValue("type of flight")),
		"number":                       strings.ToUpper(c.FormValue("number")),
		"type of aircraft":             strings.ToUpper(c.FormValue("type of aircraft")),
		"wake turbulence caat":         strings.ToUpper(c.FormValue("wake turbulence caat")),
		"equipment 1/2":                strings.ToUpper(c.FormValue("equipment 1/2")),
		"equipment 2/2":                strings.ToUpper(c.FormValue("equipment 2/2")),
		"departure":                    strings.ToUpper(c.FormValue("departure")),
		"time":                         strings.ToUpper(c.FormValue("time")),
		"cruising speed":               strings.ToUpper(c.FormValue("cruising speed")),
		"level":                        strings.ToUpper(c.FormValue("level")),
		"route1":                       strings.ToUpper(c.FormValue("route1")),
		"route2":                       strings.ToUpper(c.FormValue("route2")),
		"route3":                       strings.ToUpper(c.FormValue("route3")),
		"route4":                       strings.ToUpper(c.FormValue("route4")),
		"route5":                       strings.ToUpper(c.FormValue("route5")),
		"route6":                       strings.ToUpper(c.FormValue("route6")),
		"route7":                       strings.ToUpper(c.FormValue("route7")),
		"destination aerodrome":        strings.ToUpper(c.FormValue("destination aerodrome")),
		"total eet":                    strings.ToUpper(c.FormValue("total eet")),
		"altn aerodrome":               strings.ToUpper(c.FormValue("altn aerodrome")),
		"2nd altn aerodrome":           strings.ToUpper(c.FormValue("2nd altn aerodrome")),
		"other information1":           strings.ToUpper(c.FormValue("other information1")),
		"other information2":           strings.ToUpper(c.FormValue("other information2")),
		"other information3":           strings.ToUpper(c.FormValue("other information3")),
		"other information4":           strings.ToUpper(c.FormValue("other information4")),
		"other information5":           strings.ToUpper(c.FormValue("other information5")),
		"other information6":           strings.ToUpper(c.FormValue("other information6")),
		"other information7":           strings.ToUpper(c.FormValue("other information7")),
		"endurance":                    strings.ToUpper(c.FormValue("endurance")),
		"persons on board":             strings.ToUpper(c.FormValue("persons on board")),
		"uhf":                          strings.ToUpper(c.FormValue("uhf")),
		"vhf":                          strings.ToUpper(c.FormValue("vhf")),
		"elt":                          strings.ToUpper(c.FormValue("elt")),
		"survival":                     strings.ToUpper(c.FormValue("survival")),
		"polar":                        strings.ToUpper(c.FormValue("polar")),
		"desert":                       strings.ToUpper(c.FormValue("desert")),
		"maritime":                     strings.ToUpper(c.FormValue("maritime")),
		"jungle":                       strings.ToUpper(c.FormValue("jungle")),
		"jackets":                      strings.ToUpper(c.FormValue("jackets")),
		"light":                        strings.ToUpper(c.FormValue("light")),
		"fluores":                      strings.ToUpper(c.FormValue("fluores")),
		"uhf2":                         strings.ToUpper(c.FormValue("uhf2")),
		"vhf2":                         strings.ToUpper(c.FormValue("vhf2")),
		"dinghies":                     strings.ToUpper(c.FormValue("dinghies")),
		"number2":                      strings.ToUpper(c.FormValue("number2")),
		"capacity":                     strings.ToUpper(c.FormValue("capacity")),
		"cover":                        strings.ToUpper(c.FormValue("cover")),
		"colour":                       strings.ToUpper(c.FormValue("colour")),
		"aircraft colour and markings": strings.ToUpper(c.FormValue("aircraft colour and markings")),
		"checkbox":                     strings.ToUpper(c.FormValue("checkbox")),
		"remarks":                      strings.ToUpper(c.FormValue("remarks")),
		"pilot in command":             strings.ToUpper(c.FormValue("pilot in command")),
		"filled by":                    strings.ToUpper(c.FormValue("filled by")),
		"space reserved for additional requirements": strings.ToUpper(c.FormValue("space reserved for additional requirements")),
		"originator phone / fax 1":                   strings.ToUpper(c.FormValue("originator phone / fax 1")),
		"originator phone / fax 2":                   strings.ToUpper(c.FormValue("originator phone / fax 2")),
		"request pib 1":                              strings.ToUpper(c.FormValue("request pib 1")),
		"request pib 2":                              strings.ToUpper(c.FormValue("request pib 2")),
	}
	err := fillpdf.Fill(form, "plan.pdf", filename)
	defer os.Remove(filename)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "pdf_generation_failure",
		})
	}
	return c.File(filename)
}
