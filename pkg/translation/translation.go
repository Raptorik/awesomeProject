package translation_lesson

import (
	"awesomeProject/internal/config"
	"awesomeProject/internal/lesson"
	"awesomeProject/pkg/client/postrgresql"
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
)

var db *pgxpool.Pool

func init() {
	var err error
	db, err = postrgresql.NewClient(context.Background(), 3, config.StorageConfig{
		Username: "postgres",
		Password: "postgres",
		Host:     "localhost",
		Port:     "5432",
		Database: "postgres",
	})
	if err != nil {
		log.Fatal(err)
	}
}
func TranslateLessonName(lesson lesson.Lesson, lang string, translator Translator) error {
	rows := db.QueryRow(context.Background(), "SELECT translated_name FROM public.lesson WHERE id=$1 AND language=$2", lesson.ID, lang)
	err := rows.Scan(&lesson.TranslatedName)
	if err == nil {
		fmt.Errorf("translation already exists, no need to request from Google Cloud API")
		return nil
	} else if err != pgx.ErrNoRows {
		return fmt.Errorf("failed to query database : %v", err)
	}

	text := lesson.Name
	translatedText, err := translator.Translate(text, lang)
	if err != nil {
		return fmt.Errorf("failed to translate lesson name: %v", err)
	}

	lesson.TranslatedName = translatedText
	_, err = db.Exec(context.Background(), "INSERT INTO public.lesson (id, name, translated_name, language) VALUES ($1, $2, $3, $4)", lesson.ID, lesson.Name, lesson.TranslatedName, lang)
	if err != nil {
		log.Printf("failed to store translation in database: %v", err)
	}
	return nil
}
