package ai

import (
	"ascii-ai/config"
	"bytes"
	"context"
	"encoding/base64"
	"image"

	openai "github.com/sashabaranov/go-openai"
)

type ImageGenerator struct {
	openaiClient *openai.Client
	prompt       string
}

func NewImageGenerator(config config.ImageGeneratorConfig) *ImageGenerator {
	return &ImageGenerator{
		openai.NewClient(config.Token),
		config.Prompt,
	}
}

func (client *ImageGenerator) Generate() (image.Image, error) {
	ctx := context.Background()

	reqBase64 := openai.ImageRequest{
		Prompt:         client.prompt,
		Size:           openai.CreateImageSize1792x1024,
		ResponseFormat: openai.CreateImageResponseFormatB64JSON,
		Model:          openai.CreateImageModelDallE3,
		N:              1,
	}

	respBase64, err := client.openaiClient.CreateImage(ctx, reqBase64)
	if err != nil {
		return nil, err
	}

	imgBytes, err := base64.StdEncoding.DecodeString(respBase64.Data[0].B64JSON)
	if err != nil {
		return nil, err
	}

	r := bytes.NewReader(imgBytes)
	img, _, err := image.Decode(r)
	if err != nil {
		return nil, err
	}

	return img, nil
}
