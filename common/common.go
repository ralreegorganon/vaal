package common

type RobotState struct {
	Position     Point   `json:"position"`
	Heading      float64 `json:"heading"`
	GunHeading   float64 `json:"gunHeading"`
	RadarHeading float64 `json:"radarHeading"`
	Velocity     float64 `json:"velocity"`
	Heat         float64 `json:"heat"`
	Health       float64 `health:"health"`
	Alive        bool    `health:"health"`
}

type Point struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}
