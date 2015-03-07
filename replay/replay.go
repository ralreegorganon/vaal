package replay

type Replay struct {
	Match  string  `json:"match"`
	Arena  Arena   `json:"arena"`
	Robots []Robot `json:"robots"`
	Ticks  []Tick  `json:"ticks"`
}

type Tick struct {
	Time         int           `json:"time"`
	RobotStates  []RobotState  `json:"robotStates"`
	BulletStates []BulletState `json:"bulletStates"`
	Events       TickEvents    `json:"events"`
}

type TickEvents struct {
	ExplosionEvents []ExplosionEvent `json:"explosionEvents"`
}

type ExplosionEvent struct {
	Position Point `json:"position"`
}

type Arena struct {
	Height int   `json:"height"`
	Width  int   `json:"width"`
	Seed   int64 `json:"seed"`
}

type Robot struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type RobotState struct {
	Id           int         `json:"id"`
	Position     Point       `json:"position"`
	Heading      float64     `json:"heading"`
	GunHeading   float64     `json:"gunHeading"`
	RadarHeading float64     `json:"radarHeading"`
	Energy       float64     `json:"energy"`
	Events       RobotEvents `json:"events"`
}

type RobotEvents struct {
	RobotScannedEvents []RobotScannedEvent `json:"robotScannedEvents"`
}

type RobotScannedEvent struct {
	Distance float64 `json:"distance"`
}

type BulletState struct {
	Position Point `json:"position"`
}

type Point struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}
