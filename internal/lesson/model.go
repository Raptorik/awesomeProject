package lesson

import "awesomeProject/internal/course"

type Lesson struct {
	ID             string          `json:"id"`
	Name           string          `json:"name"`
	TranslatedName string          `json:"translated_name"`
	Language       string          `json:"language"`
	Courses        []course.Course `json:"courses"`
}
