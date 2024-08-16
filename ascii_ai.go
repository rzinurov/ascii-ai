package main

import (
	"ascii-ai/ui"
	"fmt"
	"image"
	"log"
	"os"
	"time"

	"github.com/gdamore/tcell/v2"
)

func render(s tcell.Screen, components []ui.Component, resize bool) {
	s.Clear()
	for _, c := range components {
		if resize {
			c.Resize()
		}
		c.Render()
	}
}

func initScreen() (tcell.Screen, func()) {
	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := s.Init(); err != nil {
		log.Fatalf("%+v", err)
	}
	defStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
	s.SetStyle(defStyle)
	s.Clear()

	quit := func() {
		maybePanic := recover()
		s.Fini()
		if maybePanic != nil {
			panic(maybePanic)
		}
	}

	return s, quit
}

func openImageFile(imageFilename string) (image.Image, error) {
	f, err := os.Open(imageFilename)
	if err != nil {
		return nil, err
	}

	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}

	defer f.Close()
	return img, nil
}

func main() {
	s, quit := initScreen()
	defer quit()

	components := make([]ui.Component, 0)

	imgFilename := "cat.jpg"
	img, err := openImageFile(imgFilename)
	if err != nil {
		log.Fatal("Unable to open image file: " + err.Error())
	}
	canvas := ui.NewCanvas(s, 0, 0, -1, -1)

	console := ui.NewConsole(s, 0, -5, -1, -1)
	console.Log("Welcome to ASCII-AI!")
	console.Log("Press C to toggle console")
	time.AfterFunc(1*time.Second, func() {
		canvas.Show(img)
		console.Log(fmt.Sprintf("Loaded file %s", imgFilename))
	})

	components = append(components, canvas, console)

	render(s, components, false)

	for {
		s.Show()

		ev := s.PollEvent()

		switch ev := ev.(type) {
		case *tcell.EventResize:
			render(s, components, true)
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
				return
			} else if ev.Rune() == 'c' {
				console.Toggle()
				render(s, components, false)
			}
		}
	}
}
