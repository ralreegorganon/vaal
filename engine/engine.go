package engine

import (
	"log"
	"math"
	"math/rand"

	"github.com/ralreegorganon/vaal/endpoint"
)

type Arena struct {
	Id       int
	Height   int
	Width    int
	Seed     int64
	Time     int
	Timeout  int
	Finished bool
	Robots   []*Robot
	RNG      *rand.Rand
}

func (a *Arena) Tick() {
	log.Println("Tick ", a.Time)

	for _, r := range a.Robots {
		r.Think()
	}

	// Send current state to each robot and receive commands from each robot

	// Apply commands for each robot

	// Update overall state
	a.Time += 1
	if a.Time >= a.Timeout {
		a.Finished = true
	}
}

func (a *Arena) RandomPoint() *Point {

	return nil
}

func NewArena(endpoints []*endpoint.Endpoint) *Arena {
	a := &Arena{
		Id:       0,
		Height:   800,
		Width:    800,
		Seed:     0,
		Time:     0,
		Timeout:  10,
		Finished: false,
	}

	s := rand.NewSource(a.Seed)
	a.RNG = rand.New(s)

	for _, ep := range endpoints {
		r := NewRobot(a, ep)
		a.Robots = append(a.Robots, r)
	}

	return a
}

type Robot struct {
	Id    int
	Name  string
	State RobotState
	AI    *endpoint.Endpoint
}

func NewRobot(arena *Arena, endpoint *endpoint.Endpoint) *Robot {
	p := Point{
		X: float64(arena.RNG.Intn(arena.Width)),
		Y: float64(arena.RNG.Intn(arena.Height)),
	}

	initialHeading := float64(arena.RNG.Intn(360))

	rs := RobotState{
		Position:     p,
		Heading:      initialHeading,
		GunHeading:   initialHeading,
		RadarHeading: initialHeading,
		Heat:         0,
		Velocity:     0,
		Health:       100,
		Alive:        true,
	}

	r := &Robot{
		Id:    0,
		Name:  "Test",
		State: rs,
		AI:    endpoint,
	}
	return r
}

func (r *Robot) Turn(degrees float64) {
	d := clamp(degrees, -10, 10)
	r.State.Heading += d
	r.State.GunHeading += d
	r.State.RadarHeading += d
	r.State.Heading = math.Mod(r.State.Heading, 360)
	r.State.GunHeading = math.Mod(r.State.GunHeading, 360)
	r.State.RadarHeading = math.Mod(r.State.RadarHeading, 360)
}

func (r *Robot) TurnGun(degrees float64) {
	d := clamp(degrees, -20, 20)
	r.State.GunHeading += d
	r.State.RadarHeading += d
	r.State.GunHeading = math.Mod(r.State.GunHeading, 360)
	r.State.RadarHeading = math.Mod(r.State.RadarHeading, 360)
}

func (r *Robot) TurnRadar(degrees float64) {
	d := clamp(degrees, -30, 30)
	r.State.RadarHeading += d
	r.State.RadarHeading = math.Mod(r.State.RadarHeading, 360)
}

func (r *Robot) Accelerate(velocity float64) {
	accel := clamp(velocity, -1, 1)
	r.State.Velocity += accel
	abs := math.Abs(velocity)
	r.State.Velocity = clamp(r.State.Velocity, -abs, abs)
}

func (r *Robot) Fire(energy float64) {
}

func (r *Robot) Tick() {
}

func (r *Robot) Think() {
	log.Printf("%v thinking...", r.AI.Root)
}

func clamp(val, min, max float64) float64 {
	return math.Min(math.Max(val, min), max)
}

type RobotState struct {
	Position     Point
	Heading      float64
	GunHeading   float64
	RadarHeading float64
	Velocity     float64
	Heat         float64
	Health       float64
	Alive        bool
}

type Point struct {
	X float64
	Y float64
}

type Match struct {
	Match     string
	Endpoints []*endpoint.Endpoint
}

func (m *Match) Start() {
	a := NewArena(m.Endpoints)

	for !a.Finished {
		a.Tick()
	}
}
