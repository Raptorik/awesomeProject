package translation_lesson

import (
	"cloud.google.com/go/translate"
	"context"
	"fmt"
	"golang.org/x/oauth2/google"
	"golang.org/x/text/language"
	"google.golang.org/api/option"
	"io/ioutil"
)

type Translator interface {
	Translate(text string, targetLang string) (string, error)
}

type GoogleTranslator struct {
	client *translate.Client
}

func NewGoogleTranslator(credentialsFile string) (*GoogleTranslator, error) {
	ctx := context.Background()

	creds, err := ioutil.ReadFile(credentialsFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read credentials file: %v", err)
	}

	jwtConfig, err := google.JWTConfigFromJSON(creds, translate.Scope)
	if err != nil {
		return nil, fmt.Errorf("failed to create JWT config: %v", err)
	}

	httpClient := jwtConfig.Client(ctx)

	client, err := translate.NewClient(ctx, option.WithHTTPClient(httpClient))
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
