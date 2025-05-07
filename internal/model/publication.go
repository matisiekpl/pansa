package model

type Publication struct {
	Icao string          `json:"icao"`
	Name string          `json:"name"`
	Link string          `json:"link"`
	Type PublicationType `json:"kind"`
}

type PublicationType string

const (
	PublicationTypeGEN     PublicationType = "GEN"
	PublicationTypeAD      PublicationType = "AD"
	PublicationTypeENR     PublicationType = "ENR"
	PublicationTypeUnknown PublicationType = "UNKNOWN"
)
