package ui

import (
	"github.com/gdamore/tcell/v2"
)

type Component interface {
	Resize()
	Render()
}

type BaseComponent struct {
	screen             tcell.Screen
	relativeBoundaries Boundaries
	screenBoundaries   Boundaries
	style              tcell.Style
}

type Boundaries struct {
	x1, y1, x2, y2 int
}

// Negative coordinate values will be subtracted from screen width and height at rendering time
func NewBaseComponent(screen tcell.Screen, relativeBoundaries Boundaries, style tcell.Style) *BaseComponent {
	sb := getScreenBoundaries(screen, relativeBoundaries)
	component := &BaseComponent{screen, relativeBoundaries, sb, style}
	return component
}

func getScreenBoundaries(screen tcell.Screen, rb Boundaries) Boundaries {
	xmax, ymax := screen.Size()
	return Boundaries{
		(rb.x1 + xmax) % xmax,
		(rb.y1 + ymax) % ymax,
		(rb.x2 + xmax) % xmax,
		(rb.y2 + ymax) % ymax,
	}
}

func (c *BaseComponent) getWidth() int {
	sb := c.screenBoundaries
	return sb.x2 - sb.x1 + 1
}

func (c *BaseComponent) getHeight() int {
	sb := c.screenBoundaries
	return sb.y2 - sb.y1 + 1
}

func (c *BaseComponent) clear() {
	for row := 0; row < c.getHeight(); row++ {
		for col := 0; col < c.getWidth(); col++ {
			c.drawRune(col, row, ' ')
		}
	}
}

func (c *BaseComponent) drawText(x1, y1 int, text string) {
	sb := c.screenBoundaries
	col := sb.x1 + x1
	row := sb.y1 + y1
	for _, r := range text {
		if r == '\n' {
			row++
			col = x1
			continue
		}
		if col <= sb.x2 {
			c.screen.SetContent(col, row, r, nil, c.style)
			col++
		}
		if row > sb.y2 {
			break
		}
	}
}

func (c *BaseComponent) drawRune(x, y int, r rune) {
	sb := c.screenBoundaries
	c.screen.SetContent(sb.x1+x, sb.y1+y, r, nil, c.style)
}

func (c *BaseComponent) Resize() {
	c.screenBoundaries = getScreenBoundaries(c.screen, c.relativeBoundaries)
}

func (c *BaseComponent) Render() {
	c.screen.Sync()
}
