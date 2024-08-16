package ui

import (
	"image"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/qeesung/image2ascii/convert"
)

type Canvas struct {
	*BaseComponent
	image image.Image
}

func NewCanvas(screen tcell.Screen, x1, y1, x2, y2 int) *Canvas {
	style := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
	return &Canvas{
		NewBaseComponent(screen, Boundaries{x1, y1, x2, y2}, style),
		nil,
	}
}

func (c *Canvas) Show(image image.Image) {
	c.image = image
	c.Render()
}

func (c *Canvas) Render() {
	c.clear()
	c.renderImage()
	c.BaseComponent.Render()
}

func (c *Canvas) renderImage() {
	if c.image != nil {
		content := c.generateContent()
		offset := c.getWidth()/2 - len(strings.Split(content, "\n")[0])/2
		c.drawText(offset, 0, content)
	}
}

func (c *Canvas) generateContent() string {
	convertOptions := convert.DefaultOptions
	convertOptions.FitScreen = true
	convertOptions.Colored = false

	converter := convert.NewImageConverter()
	content := converter.Image2ASCIIString(c.image, &convertOptions)

	return content
}
