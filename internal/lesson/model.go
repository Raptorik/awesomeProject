package lesson

import "awesomeProject/internal/course"

type Lesson struct {
	ID       string          `json:"id"`
	Name     string          `json:"name"`
	Language string          `json:"language"`
	Courses  []course.Course `json:"courses"`
}
