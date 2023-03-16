package student

type Student struct {
	Name     string `json:"name"`
	LessonID int    `json:"lesson_ID"`
	CourseID int    `json:"course_ID"`
}
