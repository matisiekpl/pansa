package model

type Plan struct {
	Name          string     `json:"name"`
	Distance      string     `json:"distance"`
	Duration      string     `json:"duration"`
	WindSpeed     string     `json:"windSpeed"`
	WindDirection string     `json:"windDirection"`
	Speed         string     `json:"speed"`
	Waypoints     []Waypoint `json:"waypoints"`
}
