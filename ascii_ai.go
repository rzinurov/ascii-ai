package main

import (
	"ascii-ai/ai"
	"ascii-ai/config"
	"ascii-ai/store"
	"ascii-ai/ui"
	"log"
	"sync"

	"github.com/gdamore/tcell/v2"
)

var imagePrompt string = "Cat on a synthesizer in space, high detail, realistic light, retrowave, grayscale"

func render(s tcell.Screen, components []ui.Component, resize bool) {
	s.Clear()
	for _, c := range components {
		if resize {
			c.Resize()
		}
		c.Render()
	}
}

func initScreen() tcell.Screen {
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

	return s
}

func main() {
	screen := initScreen()

	cleanup := func() {
		maybePanic := recover()
		screen.Fini()
		if maybePanic != nil {
			panic(maybePanic)
		}
	}
	defer cleanup()

	config := config.Init()

	imageGenerator := ai.NewImageGenerator(config.OpenAi.Token)

	components := make([]ui.Component, 0)
	canvas := ui.NewCanvas(screen, 0, 0, -1, -1)
	console := ui.NewConsole(screen, 0, -5, -1, -1)
	components = append(components, canvas, console)
	imageStore := store.NewImageStore(config.ImageFolder)

	nextImage := func() {
		nextImage, filename, err := imageStore.LoadNext()
		if err != nil {
			console.Log(err)
			return
		}
		canvas.Show(nextImage)
		console.Log("Loaded image ", filename)

		render(screen, components, false)
	}

	generateImageMutex := sync.Mutex{}
	generateImage := func() {
		if !generateImageMutex.TryLock() {
			console.Log("Another image generation is in progress, please wait...")
			return
		}
		defer generateImageMutex.Unlock()

		console.Log("Generating image, this may take a few seconds...")
		image, err := imageGenerator.Generate(imagePrompt)
		if err != nil {
			console.Log(err)
			return
		}
		filename, err := imageStore.Save(image)
		if err != nil {
			console.Log(err)
			return
		}
		canvas.Show(image)
		console.Log("Generated image ", filename)

		render(screen, components, false)
	}

	toggleConsole := func() {
		console.Toggle()
		render(screen, components, false)
	}

	console.Log("Welcome to ASCII-AI!")
	nextImage()
	console.Log("Press N to load next image, G to generate a new one, or C to toggle console")

	for {
		screen.Show()

		ev := screen.PollEvent()

		switch ev := ev.(type) {
		case *tcell.EventResize:
			render(screen, components, true)
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
				return
			} else if ev.Rune() == 'n' {
				go nextImage()
				render(screen, components, false)
			} else if ev.Rune() == 'g' {
				go generateImage()
			} else if ev.Rune() == 'c' {
				go toggleConsole()
			}
		}
	}
}
