package main

import (
	"log"
	"time"

	"github.com/nsf/termbox-go"
)

const (
	fgPlay        = termbox.ColorBlack
	bgPlay        = termbox.ColorBlack
	fgPlayText    = neonGreen
	bgPlayText    = termbox.ColorBlack
	flashDuration = 1 * time.Second
)

type Entity struct {
	x, y   int
	fg, bg termbox.Attribute
}

type Ball struct {
	Entity
	velX, velY int
}

var (
	paddleWidth             int
	paddleHeight            int
	leftPaddle, rightPaddle *Entity
	ball                    *Ball
)

func (g *Game) DrawPlay() {
	tbprint(ball.x, ball.y, ball.fg, ball.bg, ".")
	log.Println(paddleWidth, paddleHeight, int(float64(g.w)/100.0), int(float64(g.h)/10.0))
	tbrect(leftPaddle.x, leftPaddle.y, paddleWidth, paddleHeight, leftPaddle.fg, leftPaddle.bg, false)
	tbrect(rightPaddle.x, rightPaddle.y, paddleWidth, paddleHeight, rightPaddle.fg, rightPaddle.bg, false)
}

func (g *Game) wipePlay() {
	ball, leftPaddle, rightPaddle = nil, nil, nil
}

func (g *Game) gameOver() {
	g.GoPlay()
}

func (g *Game) UpdatePlay() {
	ball.x += ball.velX
	ball.y += ball.velY

	// check wall collision
	if ball.x < 0 {
		ball.x = 0
		ball.velX = -ball.velX
	} else if ball.x+1 > g.w {
		ball.x = g.w - 1
		ball.velX = -ball.velX
	}
	if ball.y < 0 {
		ball.y = 0
		ball.velY = -ball.velY
	} else if ball.y+1 > g.h {
		ball.y = g.h - 1
		ball.velY = -ball.velY
	}

	// check left paddle collision
	if ball.x >= (leftPaddle.x-1) && ball.x <= (leftPaddle.x+paddleWidth) &&
		ball.y >= (leftPaddle.y-1) && ball.y <= (leftPaddle.y+paddleHeight) {
		ball.x = leftPaddle.x + paddleWidth + 1
		ball.velX = -ball.velX
	}

	// check right paddle collision
	if ball.x >= (rightPaddle.x-1) && ball.x <= (rightPaddle.x+paddleWidth) &&
		ball.y >= (rightPaddle.y-1) && ball.y <= (rightPaddle.y+paddleHeight) {
		ball.x = rightPaddle.x - 1
		ball.velX = -ball.velX
	}
}

func (g *Game) HandleKeyPlay(k rune) {
	switch k {
	case 'a':
		leftPaddle.y -= 1
	case 'z':
		leftPaddle.y += 1
	case 'j':
		rightPaddle.y -= 1
	case 'm':
		rightPaddle.y += 1
	}

	if leftPaddle.y < 0 {
		leftPaddle.y = 0
	} else if (leftPaddle.y + paddleHeight) > g.h {
		leftPaddle.y = g.h - paddleHeight
	}
	if rightPaddle.y < 0 {
		rightPaddle.y = 0
	} else if (rightPaddle.y + paddleHeight) > g.h {
	}
}

func (g *Game) FreezeFlash(m string) {
	g.Draw()

	tbprint(g.w/2-len(m)/2, g.h/2, fgPlayText, bgPlayText, m)
	termbox.Flush()

	time.Sleep(flashDuration)
}

func (g *Game) GoPlay() {
	g.cfg = fgPlay
	g.cbg = bgPlay

	paddleHeight = int(float64(g.h) / 10.0)
	paddleWidth = int(float64(g.w) / 100.0)

	paddleStartY := (g.h - paddleHeight) / 2
	leftPaddle = &Entity{3, paddleStartY, white, white}
	rightPaddle = &Entity{g.w - 3 - paddleWidth, paddleStartY, white, white}
	ball = &Ball{Entity{g.w / 2, g.h / 2, white, white}, 1, 1}
}
