package models

type Replay struct {
	Id     int     `json:"id"`
	Arena  Arena   `json:"arena"`
	Robots []Robot `json:"robots"`
	Ticks  []Tick  `json:"ticks"`
}

type Tick struct {
	Time        int          `json:"time"`
	RobotStates []RobotState `json:"robotStates"`
}

type Arena struct {
	Height int `json:"height"`
	Width  int `json:"width"`
	Seed   int `json:"seed"`
}

type Robot struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type RobotState struct {
	Id           int     `json:"id"`
	X            float64 `json:"x"`
	Y            float64 `json:"y"`
	Heading      float64 `json:"heading"`
	GunHeading   float64 `json:"gunHeading"`
	RadarHeading float64 `json:"radarHeading"`
	Energy       float64 `json:"energy"`
}
