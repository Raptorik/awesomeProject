package course

import (
	"awesomeProject/internal/course"
	"awesomeProject/pkg/client/postrgresql"
	"awesomeProject/pkg/logging"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgconn"
)

type repository struct {
	client postrgresql.Client
	logger *logging.Logger
}

func (r *repository) Create(ctx context.Context, course *course.Course) error {
	q := `INSERT INTO course (name) VALUES ($1) RETURNING id`
	r.logger.Tracef(fmt.Sprintf("SQL Query: %s", q))
	if err := r.client.QueryRow(ctx, q, course.Name).Scan(&course.ID); err != nil {
		var pgErr *pgconn.PgError
		if errors.Is(err, pgErr) {
			pgErr = err.(*pgconn.PgError)
			newErr := fmt.Errorf(fmt.Sprintf("sql Error: %s, Detail: %s, Where: %s, Code: %s, SQLState: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
			r.logger.Error(newErr)
			return newErr
		}
		return err
	}
	return nil
}

func (r *repository) FindAll(ctx context.Context) (u []course.Course, err error) {
	q := `SELECT id, name FROM public.course;`
	r.logger.Tracef(fmt.Sprintf("SQL Query: %s", q))
	rows, err := r.client.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	courses := make([]course.Course, 0)

	for rows.Next() {
		var crs course.Course
		err = rows.Scan(&crs.ID, &crs.Name)
		if err != nil {
			return nil, err
		}

		courses = append(courses, crs)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return courses, nil
}

func (r *repository) FindOne(ctx context.Context, id string) (course.Course, error) {
	q := `SELECT id, name FROM public.course where id = $1;`
	r.logger.Tracef(fmt.Sprintf("SQL Query: %s", q))
	var crs course.Course
	err := r.client.QueryRow(ctx, q, id).Scan(&crs.ID, &crs.Name)
	if err != nil {
		return course.Course{}, err
	}
	return crs, nil
}

func (r *repository) Update(ctx context.Context, user course.Course) error {
	//TODO implement me
	panic("implement me")
}

func (r *repository) Delete(ctx context.Context, id string) error {
	//TODO implement me
	panic("implement me")
}

func NewRepository(client postrgresql.Client, logger *logging.Logger) course.Repository {
	return &repository{client: client, logger: logger}
}
