package model

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
