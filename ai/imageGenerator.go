package ai

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"image/png"
	"os"
	"path/filepath"
	"time"

	openai "github.com/sashabaranov/go-openai"
)

type ImageGenerator struct {
	openaiClient *openai.Client
	outputDir    string
}

func NewImageGenerator(token, outputDir string) *ImageGenerator {
	return &ImageGenerator{
		openai.NewClient(token),
		outputDir,
	}
}

func (client *ImageGenerator) Generate(prompt string) (string, error) {
	ctx := context.Background()

	reqBase64 := openai.ImageRequest{
		Prompt:         prompt,
		Size:           openai.CreateImageSize1792x1024,
		ResponseFormat: openai.CreateImageResponseFormatB64JSON,
		Model:          openai.CreateImageModelDallE3,
		N:              1,
	}

	respBase64, err := client.openaiClient.CreateImage(ctx, reqBase64)
	if err != nil {
		return "", err
	}

	imgBytes, err := base64.StdEncoding.DecodeString(respBase64.Data[0].B64JSON)
	if err != nil {
		return "", err
	}

	r := bytes.NewReader(imgBytes)
	imgData, err := png.Decode(r)
	if err != nil {
		return "", err
	}

	filename := fmt.Sprintf("%v.png", time.Now().Unix())

	file, err := os.Create(filepath.Join(client.outputDir, filename))
	if err != nil {
		return "", err
	}
	defer file.Close()

	if err := png.Encode(file, imgData); err != nil {
		return "", err
	}

	return filename, nil
}
