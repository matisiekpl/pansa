package repository

import (
	"bytes"
	_ "embed"
	"encoding/base64"
	"fmt"
	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/aymerick/raymond"
	sm "github.com/flopp/go-staticmaps"
	"github.com/golang/geo/s2"
	"github.com/matisiekpl/pansa-plan/internal/model"
	"image"
	"image/color"
	"image/png"
	"strings"
)

type ReportRepository interface {
	GenerateReport(plan model.Plan, weather []model.Weather) ([]byte, error)
}

type reportRepository struct {
}

func newReportRepository() ReportRepository {
	return &reportRepository{}
}

//go:embed assets/report.html
var reportTemplate string

func (r reportRepository) GenerateReport(plan model.Plan, weather []model.Weather) ([]byte, error) {
	mapImage, err := r.createMapImage(plan)
	var mapImageBuffer bytes.Buffer
	err = png.Encode(&mapImageBuffer, mapImage)
	if err != nil {
		return nil, err
	}
	encodedMapImage := "data:image/png;base64," + base64.StdEncoding.EncodeToString(mapImageBuffer.Bytes())
	result, err := raymond.Render(reportTemplate, map[string]interface{}{
		"plan":    plan,
		"map":     encodedMapImage,
		"weather": weather,
	})
	pdf, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		return nil, err
	}
	pdf.AddPage(wkhtmltopdf.NewPageReader(strings.NewReader(result)))
	err = pdf.Create()
	return pdf.Bytes(), nil
}

func (reportRepository) createMapImage(plan model.Plan) (image.Image, error) {
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
