package main

import (
	"ascii-ai/ai"
	"ascii-ai/config"
	"ascii-ai/ui"
	"image"
	"log"
	"os"
	"path/filepath"

	"github.com/gdamore/tcell/v2"
)

var imageDirectory string = filepath.Join(".", "images")
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

func openImageFile(imageFilename string) (image.Image, error) {
	f, err := os.Open(filepath.Join(imageDirectory, imageFilename))
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

func getImageFileNames() []string {
	files, err := os.ReadDir(imageDirectory)
	if err != nil {
		log.Fatal("Unable to read image files:", err)
	}

	fileNames := make([]string, 0)
	for _, file := range files {
		if !file.IsDir() {
			fileNames = append(fileNames, file.Name())
		}
	}
	return fileNames
}

func generateAndShowImage(imageGenerator *ai.ImageGenerator, prompt string, canvas *ui.Canvas, console *ui.Console) {
	console.Log("Generating image, this may take a few seconds...")

	filename, err := imageGenerator.Generate(prompt)

	if err != nil {
		log.Fatal("Unable to generate image: ", err)
	}

	loadAndShowImage(filename, canvas, console)
}

func loadAndShowImage(filename string, canvas *ui.Canvas, console *ui.Console) {
	image, err := openImageFile(filename)
	if err != nil {
		log.Fatal("Unable to open image file: ", err)
	}
	console.Log("Loaded image: ", filename)

	canvas.Show(image)
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

	imageGenerator := ai.NewImageGenerator(config.OpenAi.Token, imageDirectory)

	components := make([]ui.Component, 0)

	canvas := ui.NewCanvas(screen, 0, 0, -1, -1)

	console := ui.NewConsole(screen, 0, -5, -1, -1)
	console.Log("Welcome to ASCII-AI!")
	console.Log("Press G to generate a new image or C to toggle console")

	components = append(components, canvas, console)

	imageFileNames := getImageFileNames()
	lastImageFileName := imageFileNames[len(imageFileNames)-1]
	loadAndShowImage(lastImageFileName, canvas, console)

	render(screen, components, false)

	for {
		screen.Show()

		ev := screen.PollEvent()

		switch ev := ev.(type) {
		case *tcell.EventResize:
			render(screen, components, true)
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
				return
			} else if ev.Rune() == 'c' {
				console.Toggle()
				render(screen, components, false)
			} else if ev.Rune() == 'g' {
				generateAndShowImage(imageGenerator, imagePrompt, canvas, console)
				render(screen, components, false)
			}
		}
	}
}
