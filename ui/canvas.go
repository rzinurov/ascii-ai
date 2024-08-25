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

var terminalSymbolAspectRatio = 1.5

func NewCanvas(screen tcell.Screen, x1, y1, x2, y2 int, style tcell.Style) *Canvas {
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
	if c.image == nil {
		return
	}
	content := c.generateContent()
	lines := strings.Split(content, "\n")
	offsetX := c.getWidth()/2 - len(lines[0])/2
	offsetY := c.getHeight()/2 - len(lines)/2
	c.drawText(offsetX, offsetY, content)
}

func (c *Canvas) generateContent() string {
	convertOptions := convert.DefaultOptions
	convertOptions.FitScreen = false
	convertOptions.Colored = false

	imgWidth, imgHeight := c.calculateImageSize()

	convertOptions.FixedWidth = imgWidth
	convertOptions.FixedHeight = imgHeight

	converter := convert.NewImageConverter()
	content := converter.Image2ASCIIString(c.image, &convertOptions)

	return content
}

func (c *Canvas) calculateImageSize() (width, height int) {
	imgBounds := c.image.Bounds()
	imgWidth := float64(imgBounds.Size().X) * terminalSymbolAspectRatio
	imgHeight := float64(imgBounds.Size().Y)
	imgAspectRatio := imgWidth / imgHeight

	screenWidth := float64(c.getWidth())
	screenHeight := float64(c.getHeight())
	screenAspectRatio := screenWidth / screenHeight

	scaleRatio := 1.0
	if imgAspectRatio >= screenAspectRatio {
		scaleRatio = screenWidth / imgWidth
	} else {
		scaleRatio = screenHeight / imgHeight
	}

	return int(imgWidth * scaleRatio), int(imgHeight * scaleRatio)
}
