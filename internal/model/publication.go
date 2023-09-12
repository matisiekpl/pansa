package model

type Publication struct {
	Icao string `json:"icao"`
	Name string `json:"name"`
	Link string `json:"link"`
	Type string `json:"kind"`
}
