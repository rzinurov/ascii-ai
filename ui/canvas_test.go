package ui

import (
	"image"
	"image/color"
	"image/draw"
	"testing"

	"github.com/gdamore/tcell/v2"
	"github.com/stretchr/testify/assert"
)

type mockScreen struct {
	tcell.Screen
	width, height int
	content       [][]rune
}

func createMockScreen(width, height int) *mockScreen {
	c := make([][]rune, height)
	for i := range c {
		c[i] = make([]rune, width)
	}
	return &mockScreen{width: width, height: height, content: c}
}

func (ms *mockScreen) Size() (width, height int) {
	return ms.width, ms.height
}

func (ms *mockScreen) Sync() {}

func (ms *mockScreen) SetContent(x int, y int, primary rune, combining []rune, style tcell.Style) {
	ms.content[y][x] = primary
}

func (ms *mockScreen) toString() string {
	r := ""
	for i := range ms.content {
		if len(r) > 0 {
			r += "\n"
		}
		r += string(ms.content[i])
	}
	return r
}

func TestRender(t *testing.T) {
	s := tcell.StyleDefault

	imageWhite4x4 := image.NewRGBA(image.Rect(0, 0, 4, 4))
	white := color.RGBA{255, 255, 255, 255}
	draw.Draw(imageWhite4x4, imageWhite4x4.Bounds(), &image.Uniform{white}, image.Point{}, draw.Src)

	subtests := []struct {
		name                string
		x, y, width, height int
		image               image.Image
		expectedContent     string
	}{
		{
			name:            "2x2 blank",
			width:           2,
			height:          2,
			expectedContent: "  \n  ",
		},
		{
			name:            "3x4 blank",
			width:           3,
			height:          4,
			expectedContent: "   \n   \n   \n   ",
		},
		{
			name:            "3x4 blank with offset 2x1",
			x:               2,
			y:               1,
			width:           3,
			height:          4,
			expectedContent: "\x00\x00\x00\n\x00\x00 \n\x00\x00 \n\x00\x00 ",
		},
		{
			name:            "4x4 white",
			width:           4,
			height:          4,
			image:           imageWhite4x4,
			expectedContent: "    \n@@@@\n@@@@\n    ", // terminal symbols are not square
		},
	}

	for _, subtest := range subtests {
		t.Run(subtest.name, func(t *testing.T) {
			screen := createMockScreen(subtest.width, subtest.height)
			canvas := NewCanvas(screen, subtest.x, subtest.y, subtest.width-1, subtest.height-1, s)
			if subtest.image != nil {
				canvas.Show(subtest.image)
			}
			canvas.Render()

			actual := screen.toString()
			assert.Equal(t, subtest.expectedContent, actual, "Content mismatched")
		})
	}
}
