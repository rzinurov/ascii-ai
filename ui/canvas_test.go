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

	white := color.RGBA{255, 255, 255, 255}

	imageWhite4x4 := image.NewRGBA(image.Rect(0, 0, 4, 4))
	draw.Draw(imageWhite4x4, imageWhite4x4.Bounds(), &image.Uniform{white}, image.Point{}, draw.Src)

	imageWhite4x2 := image.NewRGBA(image.Rect(0, 0, 4, 2))
	draw.Draw(imageWhite4x2, imageWhite4x2.Bounds(), &image.Uniform{white}, image.Point{}, draw.Src)

	imageWhite2x4 := image.NewRGBA(image.Rect(0, 0, 2, 4))
	draw.Draw(imageWhite2x4, imageWhite2x4.Bounds(), &image.Uniform{white}, image.Point{}, draw.Src)

	imageGradient16x1 := image.NewRGBA(image.Rect(0, 0, 16, 1))
	for x := 0; x < 16; x++ {
		for y := 0; y < 1; y++ {
			draw.Draw(
				imageGradient16x1,
				image.Rectangle{Min: image.Point{x, y}, Max: image.Point{x + 1, y + 1}},
				&image.Uniform{color.RGBA{uint8(255 / 15 * x), uint8(255 / 15 * x), uint8(255 / 15 * x), 255}},
				image.Point{},
				draw.Src,
			)
		}
	}

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
			expectedContent: "    \n@@@@\n@@@@\n    ",
		},
		{
			name:            "4x4 white scaled down",
			width:           2,
			height:          2,
			image:           imageWhite4x4,
			expectedContent: "@@\n  ",
		},
		{
			name:            "4x4 white scaled up",
			width:           6,
			height:          6,
			image:           imageWhite4x4,
			expectedContent: "      \n@@@@@@\n@@@@@@\n@@@@@@\n      \n      ",
		},
		{
			name:            "4x4 white fit width",
			width:           4,
			height:          5,
			image:           imageWhite4x4,
			expectedContent: "    \n@@@@\n@@@@\n    \n    ",
		},
		{
			name:            "4x4 white fit height",
			width:           6,
			height:          4,
			image:           imageWhite4x4,
			expectedContent: "@@@@@@\n@@@@@@\n@@@@@@\n      ",
		},
		{
			name:            "2x4 white",
			width:           2,
			height:          4,
			image:           imageWhite2x4,
			expectedContent: "  \n@@\n@@\n  ",
		},
		{
			name:            "2x4 white scaled down",
			width:           2,
			height:          2,
			image:           imageWhite2x4,
			expectedContent: "@@\n@@",
		},
		{
			name:            "2x4 white scaled up",
			width:           6,
			height:          6,
			image:           imageWhite2x4,
			expectedContent: "@@@@@@\n@@@@@@\n@@@@@@\n@@@@@@\n@@@@@@\n@@@@@@",
		},
		{
			name:            "4x2 white",
			width:           4,
			height:          2,
			image:           imageWhite4x2,
			expectedContent: "@@@@\n    ",
		},
		{
			name:            "4x2 white scaled down",
			width:           2,
			height:          2,
			image:           imageWhite4x2,
			expectedContent: "@@\n  ",
		},
		{
			name:            "4x2 white scaled up",
			width:           6,
			height:          6,
			image:           imageWhite4x2,
			expectedContent: "      \n      \n@@@@@@\n      \n      \n      ",
		},
		{
			name:            "4x4 gradient",
			width:           16,
			height:          2,
			image:           imageGradient16x1,
			expectedContent: " .,:;i1tfLLCG08@\n                ",
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
