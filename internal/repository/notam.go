package repository

import (
	"encoding/json"
	"fmt"
	"github.com/matisiekpl/pansa-plan/internal/dto"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	nmsAuthorizationURL = "https://api-nms.aim.faa.gov/v1/auth/token"
	nmsNotamsURL        = "https://api-nms.aim.faa.gov/nmsapi/v1/notams"
)

type NotamRepository interface {
	Index(icao string) ([]string, error)
}

type notamRepository struct {
	config         dto.Config
	token          string
	tokenExpiresAt time.Time
	mutex          sync.Mutex
}

func newNotamRepository(config dto.Config) NotamRepository {
	return &notamRepository{config: config}
}

func (n *notamRepository) accessToken(refresh bool) (string, error) {
	n.mutex.Lock()
	defer n.mutex.Unlock()
	if !refresh && n.token != "" && time.Now().Before(n.tokenExpiresAt) {
		return n.token, nil
	}
	form := url.Values{"grant_type": {"client_credentials"}}
	request, err := http.NewRequest(http.MethodPost, nmsAuthorizationURL, strings.NewReader(form.Encode()))
	if err != nil {
		return "", err
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.SetBasicAuth(n.config.FAAKey, n.config.FAASecret)
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("nms authorization failed with status %d", response.StatusCode)
	}
	var tokenResponse dto.NmsTokenResponse
	if err := json.NewDecoder(response.Body).Decode(&tokenResponse); err != nil {
		return "", err
	}
	expiresIn, err := strconv.Atoi(tokenResponse.ExpiresIn)
	if err != nil {
		return "", err
	}
	n.token = tokenResponse.AccessToken
	n.tokenExpiresAt = time.Now().Add(time.Duration(expiresIn-60) * time.Second)
	return n.token, nil
}

func (n *notamRepository) fetch(icao string) (dto.NmsNotamResponse, error) {
	var notamResponse dto.NmsNotamResponse
	for attempt := 0; attempt < 2; attempt++ {
		token, err := n.accessToken(attempt > 0)
		if err != nil {
			return notamResponse, err
		}
		request, err := http.NewRequest(http.MethodGet, nmsNotamsURL+"?location="+url.QueryEscape(icao), nil)
		if err != nil {
			return notamResponse, err
		}
		request.Header.Set("Authorization", "Bearer "+token)
		request.Header.Set("nmsResponseFormat", "GEOJSON")
		response, err := http.DefaultClient.Do(request)
		if err != nil {
			return notamResponse, err
		}
		if response.StatusCode == http.StatusUnauthorized && attempt == 0 {
			response.Body.Close()
			continue
		}
		if response.StatusCode != http.StatusOK {
			response.Body.Close()
			return notamResponse, fmt.Errorf("nms notams request failed with status %d", response.StatusCode)
		}
		err = json.NewDecoder(response.Body).Decode(&notamResponse)
		response.Body.Close()
		return notamResponse, err
	}
	return notamResponse, fmt.Errorf("nms notams request failed with status %d", http.StatusUnauthorized)
}

func (n *notamRepository) text(feature dto.NmsNotamFeature) string {
	core := feature.Properties.CoreNOTAMData
	for _, translation := range core.NotamTranslation {
		if translation.Type == "ICAO" && translation.FormattedText != "" {
			return translation.FormattedText
		}
	}
	for _, translation := range core.NotamTranslation {
		if translation.SimpleText != "" {
			return translation.SimpleText
		}
	}
	return core.Notam.Text
}

func (n *notamRepository) Index(icao string) ([]string, error) {
	notamResponse, err := n.fetch(strings.ToUpper(icao))
	if err != nil {
		return nil, err
	}
	if notamResponse.Status != "Success" {
		var messages []string
		for _, internalError := range notamResponse.Errors {
			messages = append(messages, internalError.Message)
		}
		return nil, fmt.Errorf("nms notams request returned %s: %s", notamResponse.Status, strings.Join(messages, "; "))
	}
	notams := make([]string, 0)
	for _, feature := range notamResponse.Data.Geojson {
		text := strings.TrimSpace(strings.ReplaceAll(n.text(feature), "\r\n", "\n"))
		if text != "" {
			notams = append(notams, text)
		}
	}
	return notams, nil
}
