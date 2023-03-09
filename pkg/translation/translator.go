package translation_lesson

import (
	"context"
	"fmt"

	"cloud.google.com/go/translate"
	"golang.org/x/text/language"
	"google.golang.org/api/option"
)

type Translator interface {
	Translate(text string, targetLang string) (string, error)
}

type GoogleTranslator struct {
	client *translate.Client
}

func NewGoogleTranslator(apiKey string) (*GoogleTranslator, error) {
	ctx := context.Background()

	client, err := translate.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create translation client: %v", err)
	}

	return &GoogleTranslator{client: client}, nil
}

func (t *GoogleTranslator) Translate(text string, targetLang string) (string, error) {
	ctx := context.Background()

	target, err := language.Parse(targetLang)
	if err != nil {
		return "", fmt.Errorf("failed to parse target language: %v", err)
	}

	resp, err := t.client.Translate(ctx, []string{text}, target, nil)
	if err != nil {
		return "", fmt.Errorf("failed to translate text: %v", err)
	}

	return resp[0].Text, nil
}
