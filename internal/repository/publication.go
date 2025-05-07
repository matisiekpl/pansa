package repository

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/kaptinlin/jsonrepair"
	"github.com/matisiekpl/pansa-plan/internal/model"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strings"
	"time"
)

type PansaLanguage struct {
	ID           int    `json:"id"`
	Code         string `json:"code"`
	Flag         string `json:"flag"`
	History      string `json:"history"`
	HistoryLabel string `json:"historyLabel"`
	Help         string `json:"help"`
	HelpLabel    string `json:"helpLabel"`
	Cover        string `json:"cover"`
	CoverText    string `json:"coverText"`
	Name         string `json:"name"`
}

type PansaCommands struct {
	Title     string          `json:"title"`
	Languages []PansaLanguage `json:"languages"`
}

type PansaMenuItem struct {
	Parent    string          `json:"parent"`
	ID        string          `json:"id"`
	Href      string          `json:"href"`
	Title     string          `json:"title"`
	Depth     int             `json:"depth"`
	Level     int             `json:"level"`
	Collapsed bool            `json:"collapsed"`
	Children  []PansaMenuItem `json:"children"`
}

type PansaTabContent struct {
	Title string          `json:"title"`
	Menu  []PansaMenuItem `json:"menu"`
}

type PansaTab struct {
	ID       int                        `json:"id"`
	Title    string                     `json:"title"`
	Contents map[string]PansaTabContent `json:"contents"`
}

type PansaTabData struct {
	Generated string        `json:"generated"`
	Commands  PansaCommands `json:"commands"`
	Tabs      []PansaTab    `json:"tabs"`
}

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

func (r *publicationRepository) Index() []model.Publication {
	cacheKey := time.Now().Format(time.DateOnly)

	if cached, exists := r.cache[cacheKey]; exists {
		logrus.Infof("Found cached %d publications for %s", len(cached), cacheKey)
		return cached
	}

	link, err := r.findEAIPLink()
	if err != nil {
		logrus.Error(err)
		return nil
	}

	logrus.Infof("Found EAIP link: %s", link)

	amendmentLink, err := r.findCurrentAmendmentLink(link)
	if err != nil {
		logrus.Error(err)
		return nil
	}

	logrus.Infof("Found amendment link: %s", amendmentLink)

	tabs, err := r.extractTabs(amendmentLink)
	if err != nil {
		logrus.Error(err)
		return nil
	}

	logrus.Infof("Extracted %d tabs", len(tabs.Tabs))

	publications := r.combine(tabs, amendmentLink)
	publications = r.filterDuplicates(publications)

	logrus.Infof("Found %d publications", len(publications))

	r.cache[cacheKey] = publications
	return publications
}

func (r *publicationRepository) findEAIPLink() (string, error) {
	resp, err := http.Get("https://www.ais.pansa.pl/publikacje/aip-vfr/")
	if err != nil {
		return "", fmt.Errorf("failed to fetch URL: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("status code error: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to parse HTML: %v", err)
	}

	var eaipLink string
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		if href, exists := s.Attr("href"); exists {
			if strings.Contains(strings.ToLower(href), "eaipvfr") {
				eaipLink = href
			}
		}
	})

	if eaipLink == "" {
		return "", fmt.Errorf("no eaip link found")
	}

	return eaipLink, nil
}

func (r *publicationRepository) findCurrentAmendmentLink(eaipLink string) (string, error) {
	resp, err := http.Get(eaipLink)
	if err != nil {
		return "", fmt.Errorf("failed to fetch URL: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("status code error: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to parse HTML: %v", err)
	}

	var amendmentHref string
	foundCurrent := false

	doc.Find("h2").Each(func(i int, s *goquery.Selection) {
		if s.Text() == "Obowiązująca Zmiana" {
			foundCurrent = true
			nextTable := s.NextUntil("h2").Filter("table").First()
			nextTable.Find("a").First().Each(func(i int, a *goquery.Selection) {
				if href, exists := a.Attr("href"); exists {
					amendmentHref = href
				}
			})
		}
	})

	if !foundCurrent || amendmentHref == "" {
		return "", fmt.Errorf("no current amendment link found")
	}

	baseURL := eaipLink[:strings.LastIndex(eaipLink, "/")+1]
	return baseURL + amendmentHref, nil
}

func (r *publicationRepository) extractTabs(amendmentLink string) (*PansaTabData, error) {
	baseURL := amendmentLink[:strings.LastIndex(amendmentLink, "\\")+1]
	datasourceURL := baseURL + "v2/js/datasource.js"

	datasourceURL = strings.ReplaceAll(datasourceURL, "\\", "/")

	logrus.Info("Fetching datasource.js from:", datasourceURL)

	resp, err := http.Get(strings.ReplaceAll(datasourceURL, " ", "%20"))
	if err != nil {
		return nil, fmt.Errorf("failed to fetch datasource.js: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("status code error: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	content := string(body)
	content = strings.TrimPrefix(content, "const DATASOURCE = ")
	content = strings.TrimSuffix(content, ";")
	content = strings.ReplaceAll(content, "\t", "")

	content, err = jsonrepair.JSONRepair(content)
	if err != nil {
		return nil, fmt.Errorf("failed to repair JSON: %v", err)
	}

	var data PansaTabData
	if err := json.Unmarshal([]byte(content), &data); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %v", err)
	}

	return &data, nil
}

func (r *publicationRepository) getPDFLink(amendmentLink, href string) string {
	href = strings.ReplaceAll(href, "-pl-PL", "")
	href = strings.ReplaceAll(href, "-en-US", "")
	href = strings.ReplaceAll(href, "-en-GB", "")
	href = strings.ReplaceAll(href, ".html", ".pdf")
	baseURL := amendmentLink[:strings.LastIndex(amendmentLink, "\\")+1]
	link := baseURL + "documents/PDF/" + href + ".pdf"
	link = strings.ReplaceAll(link, " ", "%20")
	link = strings.ReplaceAll(link, "\\", "/")
	link = strings.TrimSpace(link)
	return link
}

func (r *publicationRepository) extractIcao(name string) string {
	words := strings.Fields(name)
	for _, word := range words {
		if len(word) == 4 && strings.HasPrefix(word, "EP") {
			return word
		}
	}
	return ""
}

func (r *publicationRepository) combine(tabs *PansaTabData, amendmentLink string) []model.Publication {
	var publications []model.Publication

	for _, tab := range tabs.Tabs {
		for _, content := range tab.Contents {
			var processMenuItem func(item PansaMenuItem)
			processMenuItem = func(item PansaMenuItem) {
				if item.Href != "" {
					pubType := model.PublicationTypeUnknown
					if strings.Contains(item.Href, "AD") {
						pubType = model.PublicationTypeAD
					} else if strings.Contains(item.Href, "ENR") {
						pubType = model.PublicationTypeENR
					} else if strings.Contains(item.Href, "GEN") {
						pubType = model.PublicationTypeGEN
					}

					publications = append(publications, model.Publication{
						Icao: r.extractIcao(item.Title),
						Name: strings.TrimSpace(item.Title),
						Link: r.getPDFLink(amendmentLink, item.Href),
						Type: pubType,
					})
				}

				for _, child := range item.Children {
					processMenuItem(child)
				}
			}

			for _, menu := range content.Menu {
				processMenuItem(menu)
			}
		}
	}

	filteredPublications := make([]model.Publication, 0)
	for _, publication := range publications {
		if strings.HasPrefix(publication.Link, "https") {
			filteredPublications = append(filteredPublications, publication)
		}
	}

	return filteredPublications
}

func (r *publicationRepository) filterDuplicates(publications []model.Publication) []model.Publication {
	seen := make(map[string]bool)
	var unique []model.Publication

	for _, pub := range publications {
		if !seen[pub.Link] {
			seen[pub.Link] = true
			unique = append(unique, pub)
		}
	}

	return unique
}
