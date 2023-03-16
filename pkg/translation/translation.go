package translation_lesson

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
)

func TranslateLessonName(db *pgxpool.Pool, lessonID string, lang string, translator Translator) error {
	rows := db.QueryRow(context.Background(), "SELECT name, translated_name FROM public.lesson WHERE id=$1 AND language=$2", lessonID, lang)

	var name string
	var translatedName string

	err := rows.Scan(&name, &translatedName)
	if err == nil {
		fmt.Errorf("translation already exists, no need to request from Google Cloud API")
		return nil
	} else if err != pgx.ErrNoRows {
		return fmt.Errorf("failed to query database: %v", err)
	}

	translatedText, err := translator.Translate(name, lang)
	if err != nil {
		return fmt.Errorf("failed to translate lesson name: %v", err)
	}

	translatedName = translatedText
	_, err = db.Exec(context.Background(), "INSERT INTO public.lesson (id, name, translated_name, language) VALUES ($1, $2, $3, $4)", lessonID, name, translatedName, lang)
	if err != nil {
		log.Printf("failed to store translation in database: %v", err)
	}

	return nil
}
