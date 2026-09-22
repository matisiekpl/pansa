package dto

type NmsNotamResponse struct {
	Status string             `json:"status"`
	Errors []NmsInternalError `json:"errors"`
	Data   struct {
		Geojson []NmsNotamFeature `json:"geojson"`
	} `json:"data"`
}

type NmsNotamFeature struct {
	Properties struct {
		CoreNOTAMData struct {
			Notam struct {
				Number string `json:"number"`
				Text   string `json:"text"`
			} `json:"notam"`
			NotamTranslation []struct {
				Type          string `json:"type"`
				FormattedText string `json:"formattedText"`
				SimpleText    string `json:"simpleText"`
			} `json:"notamTranslation"`
		} `json:"coreNOTAMData"`
	} `json:"properties"`
}
