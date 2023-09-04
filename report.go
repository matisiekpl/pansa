package main

import (
	"bytes"
	_ "embed"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/aymerick/raymond"
	sm "github.com/flopp/go-staticmaps"
	"github.com/golang/geo/s2"
	"github.com/labstack/echo/v4"
	"image"
	"image/color"
	"image/png"
	"net/http"
	"strings"
)

//go:embed report.html
var reportTemplate string

func generateReport(c echo.Context) error {
	var plan Plan
	err := json.Unmarshal([]byte(c.QueryParam("payload")), &plan)
	if err != nil {
		return err
	}
	mapImage, err := createMapImage(plan)
	var mapImageBuffer bytes.Buffer
	err = png.Encode(&mapImageBuffer, mapImage)
	if err != nil {
		return err
	}
	encodedMapImage := "data:image/png;base64," + base64.StdEncoding.EncodeToString(mapImageBuffer.Bytes())
	result, err := raymond.Render(reportTemplate, map[string]interface{}{
		"plan": plan,
		"map":  encodedMapImage,
	})
	pdf, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		return err
	}
	pdf.AddPage(wkhtmltopdf.NewPageReader(strings.NewReader(result)))
	err = pdf.Create()
	return c.Blob(http.StatusOK, "application/pdf", pdf.Bytes())
}

type Plan struct {
	Name      string     `json:"name"`
	Distance  string     `json:"distance"`
	Duration  string     `json:"duration"`
	Waypoints []Waypoint `json:"waypoints"`
}

type Waypoint struct {
	Name            string  `json:"name"`
	Latitude        float64 `json:"latitude"`
	Longitude       float64 `json:"longitude"`
	MagneticTrack   string  `json:"magneticTrack"`
	MagneticHeading string  `json:"magneticHeading"`
	GroundSpeed     string  `json:"groundSpeed"`
	Distance        string  `json:"distance"`
	Duration        string  `json:"duration"`
}

func createMapImage(plan Plan) (image.Image, error) {
	ctx := sm.NewContext()
	for i, waypoint := range plan.Waypoints {
		marker := sm.NewMarker(
			s2.LatLngFromDegrees(waypoint.Latitude, waypoint.Longitude),
			color.RGBA{R: 0xff, A: 0xff},
			16.0,
		)
		marker.Label = fmt.Sprintf("%d. %s", i+1, waypoint.Name)
		marker.LabelColor = color.Black
		marker.LabelYOffset = -1
		ctx.AddObject(marker)
	}
	for i := 0; i < len(plan.Waypoints)-1; i++ {
		ctx.AddObject(
			sm.NewPath(
				[]s2.LatLng{
					s2.LatLngFromDegrees(plan.Waypoints[i].Latitude, plan.Waypoints[i].Longitude),
					s2.LatLngFromDegrees(plan.Waypoints[i+1].Latitude, plan.Waypoints[i+1].Longitude),
				},
				color.RGBA{R: 0xff, A: 0xff},
				3,
			),
		)
	}
	return ctx.Render()
}
