package translation_lesson

import (
	"awesomeProject/internal/config"
	"awesomeProject/internal/lesson"
	"awesomeProject/pkg/client/postrgresql"
	"cloud.google.com/go/translate"
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/text/language"
	"google.golang.org/api/option"
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
func TranslateLessonName(lesson lesson.Lesson, lang string) error {
	rows := db.QueryRow(context.Background(), "SELECT translated_name FROM public.lesson WHERE id=$1 AND language=$3", lesson.ID, lang)
	err := rows.Scan(&lesson.TranslatedName)
	if err == nil {
		fmt.Errorf("translation already exists, no need to request from Google Cloud API")
		return nil
	} else if err != pgx.ErrNoRows {
		return fmt.Errorf("failed to query database : %v", err)
	}
	ctx := context.Background()
	tag, err := language.Parse(lang)
	if err != nil {
		return fmt.Errorf("failed to parse language tag: %v", err)
	}

	client, err := translate.NewClient(ctx, option.WithAPIKey("AIzaSyB4TPzuXJNY31cvAc4l5xEO9A3wACthmNE"))
	if err != nil {
		return fmt.Errorf("failed to create translation client: %v", err)
	}
	defer client.Close()

	resp, err := client.Translate(ctx, []string{lesson.Name}, tag, nil)
	if err != nil {
		return fmt.Errorf("failed to translate lesson name^ %v", err)
	}
	lesson.TranslatedName = resp[0].Text
	_, err = db.Exec(context.Background(), "INSERT INTO public.lesson (id, name, translated_name, language) VALUES ($1, $2, $3, $4)", lesson.ID, lesson.Name, lesson.TranslatedName, lang)
	if err != nil {
		log.Printf("failed to store translation in database: %v", err)
	}
	return nil
}
