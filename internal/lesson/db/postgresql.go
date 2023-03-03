package lesson

import (
	"awesomeProject/internal/course"
	"awesomeProject/internal/lesson"
	"awesomeProject/pkg/client/postrgresql"
	"awesomeProject/pkg/logging"
	"context"
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
func (r *repository) FindAll(ctx context.Context) (u []lesson.Lesson, err error) {
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
