package lesson

import (
	"awesomeProject/internal/course"
	"awesomeProject/internal/lesson"
	"awesomeProject/pkg/client/postrgresql"
	"awesomeProject/pkg/logging"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jackc/pgconn"
	"io/ioutil"
	"net/http"
)

type repository struct {
	client postrgresql.Client
	logger *logging.Logger
}

func NewRepository(client postrgresql.Client, logger *logging.Logger) lesson.Repository {
	return &repository{
		client: client,
		logger: logger,
	}
}
func (r *repository) FindAll(ctx context.Context) (_ []lesson.Lesson, err error) {
	q := `SELECT id, name, language FROM public.lesson;`

	rows, err := r.client.Query(ctx, q)
	if err != nil {
		return nil, err
	}

	lessons := make([]lesson.Lesson, 0)

	for rows.Next() {
		var lsn lesson.Lesson
		err = rows.Scan(&lsn.ID, &lsn.Name, &lsn.Language)
		if err != nil {
			return nil, err
		}
		sq := `SELECT course_id, name FROM student JOIN course a on a.id = student.course_id = a.id WHERE lesson_id =$1;`
		coursesRows, err := r.client.Query(ctx, sq, lsn.ID)
		if err != nil {
			return nil, err
		}

		courses := make([]course.Course, 0)

		for coursesRows.Next() {
			var crs course.Course

			err = coursesRows.Scan(&crs.ID, &crs.Name)
			if err != nil {
				return nil, err
			}
			courses = append(courses, crs)
		}
		lsn.Courses = courses
		lessons = append(lessons, lsn)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return lessons, nil
}
func (r *repository) TranslateLessonName(ctx context.Context, name string, lang string) (string, error) {
	resp, err := http.Get(fmt.Sprintf("https://translation.googleapis.com/language/translate/v2?key=%s&target=%s&q=%s", "https://awesome-project@ecstatic-baton-379809.iam.gserviceaccount.com", lang, name))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var result struct {
		Data struct {
			Translations []struct {
				TranslatedText string `json:"translatedText"`
			} `json:"translations"`
		} `json:"data"`
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", err
	}
	if len(result.Data.Translations) == 0 {
		return "", fmt.Errorf("no translations found")
	}
	return result.Data.Translations[0].TranslatedText, nil
}

func (r *repository) FindOne(ctx context.Context, id string) (lesson.Lesson, error) {
	q := `SELECT id, name, translated_name, language FROM public.lesson WHERE id = $1;`
	r.logger.Tracef(fmt.Sprintf("SQL Query: %s", q))
	var lsn lesson.Lesson
	err := r.client.QueryRow(ctx, q, id).Scan(&lsn.ID, &lsn.Name, &lsn.TranslatedName, &lsn.Language)
	if err != nil {
		return lesson.Lesson{}, err
	}
	sq := `SELECT course_id, name FROM student JOIN course a on a.id = student.course_id = a.id WHERE lesson_id = $1;`
	coursesRows, err := r.client.Query(ctx, sq, lsn.ID)
	if err != nil {
		return lesson.Lesson{}, err
	}
	courses := make([]course.Course, 0)

	for coursesRows.Next() {
		var crs course.Course
		err = coursesRows.Scan(&crs.ID, &crs.Name)
		if err != nil {
			return lesson.Lesson{}, err
		}
		courses = append(courses, crs)
	}
	lsn.Courses = courses
	return lsn, nil
}
func (r *repository) Update(ctx context.Context, lesson *lesson.Lesson) error {
	q := `UPDATE public.lesson SET translated_name = $1 WHERE id = $2;`
	r.logger.Tracef(fmt.Sprintf("SQL Query: %s", q))
	if _, err := r.client.Exec(ctx, q, lesson.TranslatedName, lesson.ID); err != nil {
		var pgErr *pgconn.PgError
		if errors.Is(err, pgErr) {
			pgErr = err.(*pgconn.PgError)
			newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
			r.logger.Error(newErr)
			return newErr
		}
		return err
	}
	return nil
}
