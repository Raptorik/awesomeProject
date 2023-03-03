package lesson

type CreateLessonDTO struct {
	Name     string `json:"name"`
	CourseID int    `json:"course_id"`
}
