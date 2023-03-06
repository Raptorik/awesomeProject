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

type repositoryCourse struct {
	client postrgresql.Client
	logger *logging.Logger
}

func NewRepositoryCourse(client postrgresql.Client, logger *logging.Logger) course.Repository {
	return &repositoryCourse{client: client, logger: logger}
}
func (r *repositoryCourse) Create(ctx context.Context, course *course.Course) error {
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

func (r *repositoryCourse) FindAll(ctx context.Context) (_ []course.Course, err error) {
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

func (r *repositoryCourse) FindOne(ctx context.Context, id string) (course.Course, error) {
	q := `SELECT id, name FROM public.course where id = $1;`
	r.logger.Tracef(fmt.Sprintf("SQL Query: %s", q))
	var crs course.Course
	err := r.client.QueryRow(ctx, q, id).Scan(&crs.ID, &crs.Name)
	if err != nil {
		return course.Course{}, err
	}
	return crs, nil
}

func (r *repositoryCourse) Update(ctx context.Context, course *course.Course) error {
	q := `UPDATE public.course SET name = $1 WHERE id = $2;`
	r.logger.Tracef(fmt.Sprintf("SQL Query: %s", q))
	if _, err := r.client.Exec(ctx, q, course.Name, course.ID); err != nil {
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

func (r *repositoryCourse) Delete(ctx context.Context, id string) error {
	q := `DELETE FROM public.course WHERE id = $1;`
	r.logger.Tracef(fmt.Sprintf("SQL Query: %s", q))
	if _, err := r.client.Exec(ctx, q, id); err != nil {
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
