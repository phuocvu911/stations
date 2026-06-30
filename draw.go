package main

// Bonus feature: the -draw flag. Renders the network as ASCII art using the
// station coordinates, then replays the schedule turn by turn with train
// positions overlaid on the map. When stdout is a terminal the frames
// animate in place; when piped they are printed one after another. The flag
// does not change the default behaviour: without -draw the program prints
// the plain schedule exactly as before.

import (
	"fmt"
	"math"
	"os"
	"strings"
	"time"
)

// extractDrawFlag removes -draw/--draw from the argument list so the
// remaining arguments can be validated as usual.
func extractDrawFlag(args []string) ([]string, bool) {
	rest := make([]string, 0, len(args))
	draw := false
	for _, a := range args {
		if a == "-draw" || a == "--draw" {
			draw = true
			continue
		}
		rest = append(rest, a)
	}
	return rest, draw
}

const (
	drawMapW   = 60 // map area, in cells
	drawMapH   = 24
	drawLabelW = 30 // room to the right of the map for station labels
	frameDelay = 800 * time.Millisecond
)

// canvas holds the static part of the picture: tracks and station markers,
// plus where each station landed after scaling. Labels and trains change
// per frame and are drawn onto a copy.
type canvas struct {
	base     [][]byte
	px, py   []int
	w, h     int
	labelAll bool // on dense maps only occupied stations get labels
}

func buildCanvas(net *Network) *canvas {
	maxX, maxY := 0, 0
	minX, minY := net.xs[0], net.ys[0]
	for i := range net.xs {
		minX, maxX = min(minX, net.xs[i]), max(maxX, net.xs[i])
		minY, maxY = min(minY, net.ys[i]), max(maxY, net.ys[i])
	}
	c := &canvas{
		px:       make([]int, len(net.xs)),
		py:       make([]int, len(net.ys)),
		labelAll: len(net.names) <= 80,
	}
	// True-to-coordinate mode, like the subject's drawing: row n is y=n,
	// the header row carries the x axis, and each x unit is wide enough
	// for its axis number. Falls back to fit-to-terminal scaling when the
	// coordinates would not fit on a screen.
	gutter := len(fmt.Sprint(maxY))
	unitW := len(fmt.Sprint(maxX)) + 1
	colOf := func(x int) int { return gutter - 1 + unitW*x }
	if maxY <= 50 && colOf(maxX) <= 78 {
		c.h = maxY + 1
		c.w = colOf(maxX) + 1 + drawLabelW
		c.clear()
		copy(c.base[0], "0")
		for x := 1; x <= maxX; x++ {
			copy(c.base[0][colOf(x):], fmt.Sprint(x))
		}
		for y := 1; y <= maxY; y++ {
			copy(c.base[y], fmt.Sprint(y))
		}
		for i := range net.xs {
			c.px[i] = colOf(net.xs[i])
			c.py[i] = net.ys[i]
		}
	} else {
		// Fit the coordinate range onto the cell grid, but cap the zoom
		// so a tiny map does not get stretched into long empty tracks.
		sx := scaleFactor(maxX-minX, drawMapW-1, 4)
		sy := scaleFactor(maxY-minY, drawMapH-1, 2)
		for i := range net.xs {
			c.px[i] = int(math.Round(float64(net.xs[i]-minX) * sx))
			c.py[i] = int(math.Round(float64(net.ys[i]-minY) * sy))
			c.w = max(c.w, c.px[i]+1)
			c.h = max(c.h, c.py[i]+1)
		}
		c.w += drawLabelW
		c.clear()
	}
	for _, conn := range net.conns {
		c.track(c.px[conn[0]], c.py[conn[0]], c.px[conn[1]], c.py[conn[1]])
	}
	for i := range net.xs {
		c.base[c.py[i]][c.px[i]] = 'X'
	}
	return c
}

func (c *canvas) clear() {
	c.base = make([][]byte, c.h)
	for y := range c.base {
		c.base[y] = []byte(strings.Repeat(" ", c.w))
	}
}

func scaleFactor(span, target int, zoomCap float64) float64 {
	if span == 0 {
		return 1
	}
	return math.Min(float64(target)/float64(span), zoomCap)
}

// track draws a line between two stations: diagonal with / or \ for as long
// as both axes have distance left, then straight with - or |. Crossings of
// different tracks become +.
func (c *canvas) track(x0, y0, x1, y1 int) {
	sx, sy := 1, 1
	if x0 > x1 {
		sx = -1
	}
	if y0 > y1 {
		sy = -1
	}
	x, y := x0, y0
	for x != x1 || y != y1 {
		var ch byte
		switch {
		case x != x1 && y != y1 && sx == sy:
			x, y, ch = x+sx, y+sy, '\\'
		case x != x1 && y != y1:
			x, y, ch = x+sx, y+sy, '/'
		case x != x1:
			x, ch = x+sx, '-'
		default:
			y, ch = y+sy, '|'
		}
		if x == x1 && y == y1 {
			break
		}
		if cur := c.base[y][x]; cur != ' ' && cur != ch {
			ch = '+'
		}
		c.base[y][x] = ch
	}
}

// renderFrame copies the base picture and writes one label per station:
// its name plus whatever trains currently occupy it.
func (c *canvas) renderFrame(net *Network, occupants []string) string {
	rows := make([][]byte, c.h)
	for y := range rows {
		rows[y] = append([]byte(nil), c.base[y]...)
	}
	for i := range net.names {
		if !c.labelAll && occupants[i] == "" {
			continue
		}
		label := net.names[i]
		if occupants[i] != "" {
			label += " " + occupants[i]
		}
		c.writeLabel(rows, c.px[i], c.py[i], label)
	}
	var b strings.Builder
	for _, row := range rows {
		b.Write([]byte(strings.TrimRight(string(row), " ")))
		b.WriteByte('\n')
	}
	return b.String()
}

// writeLabel places "X <- name" to the right of the station, or
// "name -> X" on the left when the right side is taken by another label.
func (c *canvas) writeLabel(rows [][]byte, x, y int, text string) {
	right := " <- " + text
	at := x + 1
	if at+len(right) > c.w || collides(rows[y], at+1, len(right)-1) {
		left := text + " -> "
		if start := x - len(left); start >= 0 && !collides(rows[y], start, len(left)-1) {
			copy(rows[y][start:], left)
			return
		}
	}
	// Both sides are crowded: write what fits, but never corrupt a station
	// marker or another label.
	for i := 0; i < len(right) && at+i < c.w; i++ {
		if i > 0 && occupied(rows[y][at+i]) {
			break
		}
		rows[y][at+i] = right[i]
	}
}

// collides reports whether the span already contains label text (track
// characters are fine to draw over, other labels are not).
func collides(row []byte, at, n int) bool {
	for i := at; i < at+n && i < len(row); i++ {
		if occupied(row[i]) {
			return true
		}
	}
	return false
}

func occupied(b byte) bool {
	switch b {
	case ' ', '-', '|', '/', '\\', '+':
		return false
	}
	return true
}

// drawSchedule prints the map, then one frame per turn with train positions.
func drawSchedule(net *Network, paths [][]int, counts []int) {
	c := buildCanvas(net)
	start := paths[0][0]
	end := paths[0][len(paths[0])-1]
	total := 0
	for _, n := range counts {
		total += n
	}
	animate := isTerminal()
	turn := 0
	show := func(header string, trains []simTrain) {
		occ := occupantLabels(net, paths, trains, start, end, total)
		if animate && turn > 0 {
			fmt.Print("\x1b[H\x1b[2J")
		}
		fmt.Println()
		fmt.Println(header)

		fmt.Print(c.renderFrame(net, occ))
		fmt.Println()
		if animate {
			time.Sleep(frameDelay)
		}
	}
	show(fmt.Sprintf("Moving %d trains, %s -> %s", total, net.names[start], net.names[end]), nil)
	simulate(paths, counts, net.names, func(movements []string, trains []simTrain) {
		turn++
		show(fmt.Sprintf("Turn %d: %s", turn, strings.Join(movements, " ")), trains)
	})
}

// occupantLabels builds the per-station suffix shown next to each name:
// trains sitting at the station, or waiting/arrived tallies at the ends.
func occupantLabels(net *Network, paths [][]int, trains []simTrain, start, end, total int) []string {
	occ := make([]string, len(net.names))
	arrived := 0
	perStation := make(map[int][]string)
	for i, t := range trains {
		if t.done {
			arrived++
			continue
		}
		at := paths[t.path][t.pos]
		perStation[at] = append(perStation[at], fmt.Sprintf("T%d", i+1))
	}
	for at, names := range perStation {
		if len(names) > 3 {
			names = append(names[:3], fmt.Sprintf("+%d", len(names)-3))
		}
		occ[at] = "[" + strings.Join(names, ",") + "]"
	}
	if waiting := total - len(trains); waiting > 0 {
		occ[start] = fmt.Sprintf("[%d waiting]", waiting)
	}
	if arrived > 0 {
		occ[end] = fmt.Sprintf("[%d arrived]", arrived)
	}
	return occ
}

func isTerminal() bool {
	fi, err := os.Stdout.Stat()
	return err == nil && fi.Mode()&os.ModeCharDevice != 0
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
