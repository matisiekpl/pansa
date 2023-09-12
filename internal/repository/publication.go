package repository

import (
	"github.com/anaskhan96/soup"
	"github.com/matisiekpl/pansa-plan/internal/model"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
	"time"
)

type PublicationRepository interface {
	Index() []model.Publication
}

type publicationRepository struct {
	cache map[string][]model.Publication
}

func newPublicationRepository() PublicationRepository {
	return &publicationRepository{
		cache: make(map[string][]model.Publication),
	}
}

func (publicationRepository) query(source string) []model.Publication {
	var publications []model.Publication
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
				publications = append(publications, model.Publication{
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
				publications = append(publications, model.Publication{
					Icao: "",
					Name: name,
					Link: url,
					Type: "SUP",
				})
			}
			if strings.Contains(url, "_GEN_") {
				publications = append(publications, model.Publication{
					Icao: "",
					Name: name,
					Link: url,
					Type: "GEN",
				})
			}
			if strings.Contains(url, "_ENR_") {
				publications = append(publications, model.Publication{
					Icao: "",
					Name: name,
					Link: url,
					Type: "ENR",
				})
			}
		}
	}
	return publications
}

func (p publicationRepository) Index() []model.Publication {
	p.query("https://www.ais.pansa.pl/publikacje/aip-vfr/")
	if val, ok := p.cache[time.Now().Format(time.DateOnly)]; ok {
		return val
	}
	p.cache[time.Now().Format(time.DateOnly)] = append(p.query("https://www.ais.pansa.pl/publikacje/aip-vfr/"), p.query("https://www.ais.pansa.pl/publikacje/aip-ifr/")...)
	logrus.Infof("Downloaded %d AIPs", len(p.cache[time.Now().Format(time.DateOnly)]))
	return p.cache[time.Now().Format(time.DateOnly)]
}
