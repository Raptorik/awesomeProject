package lesson

import (
	"awesomeProject/internal/entities/course"
)

type Lesson struct {
	ID             string          `json:"id"`
	Name           string          `json:"name"`
	TranslatedName string          `json:"translated_name"`
	Courses        []course.Course `json:"courses"`
}

var l = Lesson{
	ID:             `1`,
	Name:           "Курс Голанг легкий",
	TranslatedName: "Golang beginner course",
	Courses:        []course.Course{},
}
