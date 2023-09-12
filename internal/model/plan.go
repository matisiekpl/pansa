package model

type Plan struct {
	Name      string     `json:"name"`
	Distance  string     `json:"distance"`
	Duration  string     `json:"duration"`
	Waypoints []Waypoint `json:"waypoints"`
}
