package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/nsf/termbox-go"
)

func tbprint(x, y int, fg, bg termbox.Attribute, msg string) {
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x++
	}
}

func tbrect(x, y, w, h int, fg, bg termbox.Attribute, border bool) {
	end := " " + strings.Repeat("_", w)
	if border {
		tbprint(x, y-1, fg, bg, end)
	}

	s := strings.Repeat(" ", w)
	if border {
		s = fmt.Sprintf("%c%s%c", '|', s, '|')
	}

	for i := 0; i < h; i++ {
		tbprint(x, y, fg, bg, s)
		y++
	}

	if border {
		tbprint(x, y, fg, bg, end)
	}
}

const (
	fgDefault = termbox.ColorRed
	bgDefault = termbox.ColorYellow
	fps       = 30
)

type Game struct {
	evq   chan termbox.Event
	timer <-chan time.Time

	// frame counter
	fc uint8
	w  int
	h  int

	// fg and bg colors used when termbox.Clear() is called
	cfg termbox.Attribute
	cbg termbox.Attribute
}

func NewGame() *Game {
	return &Game{
		evq:   make(chan termbox.Event),
		timer: time.Tick(time.Duration(1000/fps) * time.Millisecond),
		fc:    1,
	}
}

// Tick allows us to rate limit the FPS
func (g *Game) Tick() {
	<-g.timer
	g.fc++
	if g.fc > fps {
		g.fc = 1
	}
}

func (g *Game) Listen() {
	go func() {
		for {
			g.evq <- termbox.PollEvent()
		}
	}()
}

func (g *Game) HandleKey(k rune) {
	g.HandleKeyPlay(k)
}

func (g *Game) FitScreen() {
	termbox.Clear(g.cfg, g.cbg)
	g.w, g.h = termbox.Size()
	//g.Draw()
}

func (g *Game) Draw() {
	termbox.Clear(g.cfg, g.cbg)
	g.DrawPlay()
	termbox.Flush()
}

func (g *Game) Update() {
	g.Tick()
	g.UpdatePlay()
	return
}

func main() {
	if err := termbox.Init(); err != nil {
		log.Fatalln(err)
	}
	termbox.SetOutputMode(termbox.Output256)
	defer termbox.Close()

	f, err := os.Create(".l")
	if err != nil {
		log.Fatalln(err)
	}
	log.SetOutput(f)

	g := NewGame()

	g.Listen()
	g.FitScreen()
	g.GoPlay()

main:
	for {
		select {
		case ev := <-g.evq:
			switch ev.Type {
			case termbox.EventKey:
				if ev.Ch != 0 {
					if ev.Ch == 'q' {
						break main
					}
					g.HandleKey(ev.Ch)
				}
			case termbox.EventResize:
				g.FitScreen()
			}
		default:
		}

		g.Update()
		g.Draw()
	}
}
