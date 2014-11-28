package engine

import (
	"log"
	"math"
	"math/rand"

	"github.com/ralreegorganon/vaal/common"
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
		if !r.State.Alive {
			continue
		}
		r.Scan()
	}

	// Send current state to each robot and receive commands from each robot
	toProcess := make(map[*Robot]*common.RobotCommands, 0)
	for _, r := range a.Robots {
		if !r.State.Alive {
			continue
		}
		commands := r.Think()
		if commands != nil {
			toProcess[r] = commands
		}
	}

	// Apply commands for each robot
	for r, c := range toProcess {
		if !r.State.Alive {
			continue
		}
		r.Tick(c)
	}

	// Constrain
	for _, r := range a.Robots {
		r.State.Position.X = clamp(r.State.Position.X, 0, float64(a.Width))
		r.State.Position.Y = clamp(r.State.Position.Y, 0, float64(a.Height))
	}

	// Update overall state
	a.Time += 1
	if a.Time >= a.Timeout {
		a.Finished = true
	}

	alive := 0
	for _, r := range a.Robots {
		if r.State.Alive {
			alive++
		}
	}

	if alive <= 1 {
		a.Finished = true
	}
}

func (a *Arena) RandomPoint() *common.Point {

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
	State *common.RobotState
	AI    *endpoint.Endpoint
}

func NewRobot(arena *Arena, endpoint *endpoint.Endpoint) *Robot {
	p := common.Point{
		X: float64(arena.RNG.Intn(arena.Width)),
		Y: float64(arena.RNG.Intn(arena.Height)),
	}

	initialHeading := float64(arena.RNG.Intn(360))

	rs := &common.RobotState{
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

func (r *Robot) Tick(commands *common.RobotCommands) {
	r.Scan()

	if commands.Fire != 0 {
		r.Fire(commands.Fire)
	}
	if commands.Turn != 0 {
		r.Turn(commands.Turn)
	}
	if commands.TurnGun != 0 {
		r.TurnGun(commands.TurnGun)
	}
	if commands.TurnRadar != 0 {
		r.TurnRadar(commands.TurnRadar)
	}
	if commands.Accelerate != 0 {
		r.Accelerate(commands.Accelerate)
	}

	r.Cool()

	r.State.Position.X += math.Sin(r.State.Heading*math.Pi/180) * r.State.Velocity
	r.State.Position.Y += math.Cos(r.State.Heading*math.Pi/180) * r.State.Velocity
}

func (r *Robot) Scan() {
}

func (r *Robot) Cool() {
	cooled := r.State.Heat - 0.1
	r.State.Heat = clamp(cooled, 0, cooled)
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

func (r *Robot) Think() *common.RobotCommands {
	log.Printf("%v thinking...", r.AI.Root)
	log.Printf("%+v\n", r.State)
	commands, err := r.AI.Think(r.State)
	if err != nil {
		log.Printf("%v error from %v, killing it", err, r.AI.Root)
		r.State.Alive = false
		return nil
	}
	return commands
}

func clamp(val, min, max float64) float64 {
	return math.Max(math.Min(val, max), min)
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
