package main

import (
	"github.com/anaskhan96/soup"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"strings"
	"time"
)

type AIP struct {
	Icao string `json:"icao"`
	Name string `json:"name"`
	Link string `json:"link"`
	Type string `json:"kind"`
}

var cache = make(map[string][]AIP)

func getADs() []AIP {
	if val, ok := cache[time.Now().Format(time.DateOnly)]; ok {
		return val
	}
	cache[time.Now().Format(time.DateOnly)] = append(listADs("https://www.ais.pansa.pl/publikacje/aip-vfr/"), listADs("https://www.ais.pansa.pl/publikacje/aip-ifr/")...)
	logrus.Infof("Downloaded %d AIPs", len(cache[time.Now().Format(time.DateOnly)]))
	return cache[time.Now().Format(time.DateOnly)]
}

func listADs(source string) []AIP {
	var aips []AIP
	resp, err := soup.Get(source)
	if err != nil {
		os.Exit(1)
	}
	doc := soup.HTMLParse(resp)
	elements := doc.FindAll("details")
	for _, element := range elements {
		if element.Attrs()["class"] == "sub-details" {
			label := element.Children()[0].Text()
			icao := strings.TrimSpace(strings.ReplaceAll(label, "VFR", ""))
			for _, link := range element.Children()[1].FindAll("a") {
				url := strings.ReplaceAll(strings.ReplaceAll(link.Attrs()["href"], " ", ""), "	", "")
				name := strings.TrimSpace(strings.ReplaceAll(link.Text(), "\t", ""))
				aips = append(aips, AIP{
					Icao: icao,
					Name: name,
					Link: url,
					Type: "AD",
				})
			}
		}

		for _, link := range element.FindAll("a") {
			url := strings.ReplaceAll(strings.ReplaceAll(link.Attrs()["href"], " ", ""), "	", "")
			name := strings.TrimSpace(strings.ReplaceAll(link.Text(), "\t", ""))
			if strings.Contains(url, "_Sup_") {
				aips = append(aips, AIP{
					Icao: "",
					Name: name,
					Link: url,
					Type: "SUP",
				})
			}
			if strings.Contains(url, "_GEN_") {
				aips = append(aips, AIP{
					Icao: "",
					Name: name,
					Link: url,
					Type: "GEN",
				})
			}
			if strings.Contains(url, "_ENR_") {
				aips = append(aips, AIP{
					Icao: "",
					Name: name,
					Link: url,
					Type: "ENR",
				})
			}
		}
	}
	return aips
}

func serveAD(c echo.Context) error {
	return c.JSON(http.StatusOK, getADs())
}
