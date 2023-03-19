package repository

import (
	"awesomeProject/internal/entities/course"
	lesson2 "awesomeProject/internal/entities/lesson"
	"awesomeProject/pkg/client/postgresql"
	"awesomeProject/pkg/logging"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgconn"
)

type repositoryLesson struct {
	client postgresql.Client
	logger *logging.Logger
}

func NewRepositoryLesson(client postgresql.Client, logger *logging.Logger) lesson2.Repository {
	return &repositoryLesson{client: client, logger: logger}
}
func (r *repositoryLesson) FindAll(ctx context.Context) (_ []lesson2.Lesson, err error) {
	q := `SELECT id, name, language FROM public.lesson;`

	rows, err := r.client.Query(ctx, q)
	if err != nil {
		return nil, err
	}

	lessons := make([]lesson2.Lesson, 0)

	for rows.Next() {
		var lsn lesson2.Lesson
		err = rows.Scan(&lsn.ID, &lsn.Name)
		if err != nil {
			return nil, err
		}
		sq := `SELECT course_id, name FROM student JOIN course a on a.id = student.course_id AND a.id = student.course_id WHERE lesson_id =$1;`
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

func (r *repositoryLesson) FindOne(ctx context.Context, id string) (lesson2.Lesson, error) {
	q := `SELECT id, name, translated_name, language FROM public.lesson WHERE id = $1;`
	r.logger.Tracef(fmt.Sprintf("SQL Query: %s", q))
	var lsn lesson2.Lesson
	err := r.client.QueryRow(ctx, q, id).Scan(&lsn.ID, &lsn.Name, &lsn.TranslatedName)
	if err != nil {
		return lesson2.Lesson{}, err
	}
	sq := `SELECT course_id, name FROM student JOIN course a ON student.course_id = a.id WHERE lesson_id = $1;`
	coursesRows, err := r.client.Query(ctx, sq, lsn.ID)
	if err != nil {
		return lesson2.Lesson{}, err
	}
	courses := make([]course.Course, 0)

	for coursesRows.Next() {
		var crs course.Course
		err = coursesRows.Scan(&crs.ID, &crs.Name)
		if err != nil {
			return lesson2.Lesson{}, err
		}
		courses = append(courses, crs)
	}
	lsn.Courses = courses
	return lsn, nil
}
func (r *repositoryLesson) Update(ctx context.Context, lesson *lesson2.Lesson) error {
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
