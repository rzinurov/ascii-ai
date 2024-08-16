package ui

import (
	"math"

	"github.com/gdamore/tcell/v2"
)

type Console struct {
	*BaseComponent
	messages  []string
	nextIndex int
	visible   bool
}

func NewConsole(screen tcell.Screen, x1, y1, x2, y2 int) *Console {
	numRows := int(math.Abs(float64(y2 - y1 - 1)))
	style := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
	return &Console{
		NewBaseComponent(screen, Boundaries{x1, y1, x2, y2}, style),
		make([]string, numRows),
		0,
		true,
	}
}

func (c *Console) Log(text string) {
	if c.nextIndex > len(c.messages)-1 {
		c.shiftMessages()
		c.nextIndex = len(c.messages) - 1
	}
	c.messages[c.nextIndex] = text
	c.nextIndex++
	c.Render()
}

func (c *Console) shiftMessages() {
	for i := 0; i < len(c.messages)-1; i++ {
		c.messages[i] = c.messages[i+1]
	}
}

func (c *Console) Render() {
	if !c.visible {
		return
	}
	c.clear()
	c.renderBackground()
	c.renderMessages()
	c.BaseComponent.Render()
}

func (c *Console) renderBackground() {
	width, height := c.getWidth(), c.getHeight()

	// Draw borders
	for col := 1; col < width-1; col++ {
		c.drawRune(col, 0, tcell.RuneHLine)
		c.drawRune(col, height-1, tcell.RuneHLine)
	}
	for row := 1; row < height-1; row++ {
		c.drawRune(0, row, tcell.RuneVLine)
		c.drawRune(width-1, row, tcell.RuneVLine)
	}

	// // Draw corners
	c.drawRune(0, 0, tcell.RuneULCorner)
	c.drawRune(width-1, 0, tcell.RuneURCorner)
	c.drawRune(0, height-1, tcell.RuneLLCorner)
	c.drawRune(width-1, height-1, tcell.RuneLRCorner)

	// Draw title
	c.drawText(1, 0, " Console ")
}

func (c *Console) renderMessages() {
	for i := 0; i < len(c.messages); i++ {
		c.drawText(1, 1+i, c.messages[i])
	}
}

func (c *Console) Toggle() {
	c.visible = !c.visible
}
