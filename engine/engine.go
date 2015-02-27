package engine

import (
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/ralreegorganon/vaal/common"
	"github.com/ralreegorganon/vaal/endpoint"
	"github.com/ralreegorganon/vaal/replay"
)

type Arena struct {
	Match    string
	Height   int
	Width    int
	Seed     int64
	Time     int
	Timeout  int
	Finished bool
	Robots   []*Robot
	Bullets  []*Bullet
	RNG      *rand.Rand
}

func NewArena(match string, endpoints []*endpoint.Endpoint) *Arena {
	a := &Arena{
		Match:    match,
		Height:   800,
		Width:    800,
		Seed:     time.Now().Unix(),
		Time:     0,
		Timeout:  100000,
		Finished: false,
		Robots:   make([]*Robot, 0),
		Bullets:  make([]*Bullet, 0),
	}

	s := rand.NewSource(a.Seed)
	a.RNG = rand.New(s)

	for i, ep := range endpoints {
		r := NewRobot(a, ep)
		r.Id = i
		a.Robots = append(a.Robots, r)
	}

	return a
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
		r.Tick(c, a)
	}

	// Constrain
	for _, r := range a.Robots {
		r.State.Position.X = clamp(r.State.Position.X, 0, float64(a.Width))
		r.State.Position.Y = clamp(r.State.Position.Y, 0, float64(a.Height))
	}

	// Tick bullets
	liveBullets := make([]*Bullet, 0)
	for _, b := range a.Bullets {
		b.Tick(a)
		if b.Alive {
			liveBullets = append(liveBullets, b)
		}
	}

	a.Bullets = liveBullets

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

func (a *Arena) RandomPoint() common.Point {
	p := common.Point{
		X: float64(a.RNG.Intn(a.Width)),
		Y: float64(a.RNG.Intn(a.Height)),
	}

	return p
}

func (a *Arena) AddBullet(bullet *Bullet) {
	a.Bullets = append(a.Bullets, bullet)
}

type Robot struct {
	Id    int
	Name  string
	State *common.RobotState
	AI    *endpoint.Endpoint
}

func NewRobot(arena *Arena, endpoint *endpoint.Endpoint) *Robot {
	p := arena.RandomPoint()

	initialHeading := float64(arena.RNG.Intn(360))

	rs := &common.RobotState{
		Position:     p,
		Heading:      0,
		GunHeading:   initialHeading,
		RadarHeading: initialHeading,
		Heat:         0,
		Velocity:     0,
		Energy:       100,
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

func (r *Robot) Tick(commands *common.RobotCommands, arena *Arena) {
	r.Scan()

	if commands.Fire != 0 {
		bullet := r.Fire(commands.Fire)
		if bullet != nil {
			arena.AddBullet(bullet)
		}
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
	r.State.Position.Y += -math.Cos(r.State.Heading*math.Pi/180) * r.State.Velocity
}

func (r *Robot) Hit(bullet *Bullet) {
	r.State.Energy -= bullet.Firepower
	if r.State.Energy <= 0 {
		r.State.Energy = 0
		r.State.Alive = false
	}
	// damage taken event
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

func (r *Robot) Fire(energy float64) *Bullet {
	if r.State.Heat != 0 {
		return nil
	}

	fp := clamp(energy, 0.1, 3)
	r.State.Heat += 1.0 + fp/5.0

	bullet := &Bullet{
		Firepower: fp * 4,
		Alive:     true,
		Heading:   r.State.GunHeading,
		Position:  r.State.Position,
		Velocity:  20 - 3*fp,
		Origin:    r,
	}

	return bullet
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
	Replay    *replay.Replay
}

func (m *Match) Start() {
	a := NewArena(m.Match, m.Endpoints)

	m.SetupReplayForArena(a)

	m.UpdateReplayForTick(a)
	for !a.Finished {
		a.Tick()
		m.UpdateReplayForTick(a)
	}
	log.Println("------------------------------------match done")
}

func (m *Match) SetupReplayForArena(arena *Arena) {
	m.Replay = &replay.Replay{
		Match: arena.Match,
		Arena: replay.Arena{
			Width:  arena.Width,
			Height: arena.Height,
			Seed:   arena.Seed,
		},
		Robots: make([]replay.Robot, 0),
		Ticks:  make([]replay.Tick, 0),
	}

	for _, r := range arena.Robots {
		log.Printf("%+v\n", r)
		b := replay.Robot{
			Id:   r.Id,
			Name: r.Name,
		}
		m.Replay.Robots = append(m.Replay.Robots, b)
	}
}

func (m *Match) UpdateReplayForTick(arena *Arena) {
	t := replay.Tick{
		Time:         arena.Time,
		RobotStates:  make([]replay.RobotState, 0),
		BulletStates: make([]replay.BulletState, 0),
	}

	for _, r := range arena.Robots {
		s := replay.RobotState{
			Id: r.Id,
			Position: replay.Point{
				X: r.State.Position.X,
				Y: r.State.Position.Y,
			},
			Heading:      r.State.Heading,
			GunHeading:   r.State.GunHeading,
			RadarHeading: r.State.RadarHeading,
			Energy:       r.State.Energy,
		}
		t.RobotStates = append(t.RobotStates, s)
	}

	for _, b := range arena.Bullets {
		s := replay.BulletState{
			Position: replay.Point{
				X: b.Position.X,
				Y: b.Position.Y,
			},
		}
		t.BulletStates = append(t.BulletStates, s)
	}

	m.Replay.Ticks = append(m.Replay.Ticks, t)
}

type Bullet struct {
	Firepower float64
	Heading   float64
	Velocity  float64
	Position  common.Point
	Alive     bool
	Origin    *Robot
}

func (b *Bullet) Tick(arena *Arena) {
	b.Position.X += math.Sin(b.Heading*math.Pi/180) * b.Velocity
	b.Position.Y += -math.Cos(b.Heading*math.Pi/180) * b.Velocity

	x_max := float64(arena.Width)
	y_max := float64(arena.Height)

	if b.Position.X < 0 || b.Position.Y < 0 || b.Position.X > x_max || b.Position.Y > y_max {
		b.Alive = false
	}

	for _, r := range arena.Robots {
		if r != b.Origin && r.State.Alive {
			if math.Hypot(b.Position.Y-r.State.Position.Y, r.State.Position.X-b.Position.X) < 40 {
				b.Alive = false
				r.Hit(b)
			}
		}
	}
}
