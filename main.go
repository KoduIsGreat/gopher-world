package main

import (
	tl "github.com/JoelOtter/termloop"
	"math/rand"
)

// state needed by most entities in the simulation
type state struct {
	flock     *tl.Entity
	tickSpeed float64
}
// The player Implements Tick for handling user input
// can pan around the world view by holding middle mouse
// can left click to set a location for goblins to flock to
// can right click to dismiss flock command
type player struct {
	*tl.Entity
	panX int
	panY int
	game *state
	panning bool
	level *tl.BaseLevel
}
func (c *player) Tick(event tl.Event) {
	if event.Type == tl.EventMouse { // Is it a mouse event?
		x, y := event.MouseX, event.MouseY
		offsetX, offsetY :=c.level.Offset()
		if c.panning {
			offsetX -= x - c.panX
			offsetY -= y - c.panY
			c.panX = x
			c.panY = y
			c.level.SetOffset(offsetX, offsetY)
		}
		switch event.Key { // If so, switch on the pressed mouse button.
		case tl.MouseMiddle:
			c.panning = true
			c.panX, c.panY = x, y
		case tl.MouseRelease:
			c.panning =false
			c.panX, c.panY = 0, 0
		case tl.MouseLeft:
			c.level.RemoveEntity(c.game.flock)
			c.game.flock = tl.NewEntity(x-offsetX,y-offsetY,1,1)
			c.game.flock.SetCell(0,0, &tl.Cell{Bg: tl.ColorWhite, Fg: tl.ColorRed, Ch:'!'})
			c.level.AddEntity(c.game.flock)
		case tl.MouseRight:
			c.level.RemoveEntity(c.game.flock)
			c.game.flock = nil
		}
	}
}
// directionality
type direction int
var (
	gopherColor = tl.ColorCyan
)
const (
	UP direction = iota
	DOWN
	LEFT
	RIGHT
)

// common attributes
type attributes struct {
	health int
	dmg int
	mspd int
	aspd int
	fly bool
}
// Our character, apply logic in the Draw function
type gopher struct {
	*tl.Entity
	game *state
	interval float64
	facing direction
	prevX, prevY int
	attributes
}
func (g *gopher) Draw(screen *tl.Screen) {
	g.interval += screen.TimeDelta()
	if g.interval > g.game.tickSpeed {
		g.prevX, g.prevY = g.Position()
		dir := direction(rand.Intn(3))
		if g.game.flock != nil {
			x, y := g.game.flock.Position()
			if g.prevX < x {
				dir = RIGHT
			} else if g.prevX > x {
				dir = LEFT
			} else if g.prevY > y {
				dir = DOWN
			} else if g.prevY < y {
				dir = UP
			}
		}

		switch dir {
		case RIGHT:
			gopherRight(g, tl.Cell{Bg: gopherColor})
			g.SetPosition(g.prevX+g.attributes.mspd, g.prevY)
		case LEFT:
			gopherLeft(g,tl.Cell{Bg: gopherColor})
			g.SetPosition(g.prevX-g.attributes.mspd, g.prevY)
		case UP:
			gopherDown(g, tl.Cell{Bg: gopherColor})
			g.SetPosition(g.prevX, g.prevY+g.attributes.mspd)
		case DOWN:
			gopherUp(g, tl.Cell{Bg: gopherColor})
			g.SetPosition(g.prevX, g.prevY-g.attributes.mspd)
		}
		g.interval = 0
	}
	g.Entity.Draw(screen)
}


func (g *gopher) Collide(collision tl.Physical) {
	// Check if it's a Rectangle we're colliding with
	if _, ok := collision.(*tl.Rectangle); ok {
		g.SetPosition(g.prevX, g.prevY)
	} else if _, ok := collision.(*gopher); ok {
		g.SetPosition(g.prevX, g.prevY)
	}
}

func gopherLeft(g *gopher, cell tl.Cell) {
	canvas := tl.NewCanvas(2, 2)
	canvas[1][0] = cell
	canvas[1][1] = cell
	canvas[0][0] = tl.Cell{Bg: gopherColor, Fg: tl.ColorWhite, Ch:'o'}
	g.SetCanvas(&canvas)
}

func gopherRight(g *gopher, cell tl.Cell) {
	canvas := tl.NewCanvas(2, 2)
	canvas[0][0] = cell
	canvas[1][0] = tl.Cell{Bg: gopherColor, Fg: tl.ColorWhite, Ch:'o'}
	canvas[0][1] = cell
	g.SetCanvas(&canvas)
}

func gopherUp(g *gopher, cell tl.Cell) {
	canvas := tl.NewCanvas(3, 2)
	canvas[0][0] = cell
	canvas[1][0] = tl.Cell{Bg: gopherColor, Fg: tl.ColorWhite, Ch:'#'}
	canvas[2][0] = cell
	canvas[1][1] = cell
	g.SetCanvas(&canvas)
}

func gopherDown(g *gopher, cell tl.Cell) {
	canvas := tl.NewCanvas(3, 2)
	canvas[0][0] = tl.Cell{Bg: gopherColor, Fg: tl.ColorWhite, Ch:'o'}
	canvas[1][0] = tl.Cell{Bg: gopherColor, Fg: tl.ColorBlack, Ch:'_'}
	canvas[2][0] = tl.Cell{Bg: gopherColor, Fg: tl.ColorWhite, Ch:'o'}
	canvas[1][1] = cell
	g.SetCanvas(&canvas)
}

func newGopher(x, y int, state *state, cell tl.Cell) *gopher {
	canvas := tl.NewCanvas(2, 2)
	canvas[0][0] = cell
	canvas[1][0] = tl.Cell{Bg: gopherColor, Fg: tl.ColorWhite, Ch:'o'}
	canvas[0][1] = cell
	return &gopher{
		Entity: tl.NewEntityFromCanvas(x,y, canvas),
		game: state,
		attributes: attributes{
			health: 40,
			dmg:    5,
			mspd:    1,
			aspd: 1, // attacks per second
		},
	}
}



func main() {
	game := tl.NewGame()
	game.Screen().SetFps(60)
	level := tl.NewBaseLevel(tl.Cell{
			Bg: tl.ColorGreen,
			Fg: tl.ColorBlack,
		})
	state := &state{tickSpeed: .8}
	camera := &player{
		Entity: tl.NewEntity(1,1,0,0),
		game: state,
		level: level,
	}
	level.AddEntity(tl.NewRectangle(10, 10, 50, 20, tl.ColorBlue))
	level.AddEntity(camera)
	goblin := newGopher(4,4, state, tl.Cell{Bg: gopherColor})
	goblin2 := newGopher(4, 8,state, tl.Cell{Bg: gopherColor})
	goblin3 := newGopher(4, 12,state, tl.Cell{Bg: gopherColor})
	level.AddEntity(goblin3)
	level.AddEntity(goblin2)
	level.AddEntity(goblin)
	game.Screen().SetLevel(level)
	game.Start()
}