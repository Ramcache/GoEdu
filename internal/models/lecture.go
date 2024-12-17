package models

type Lecture struct {
	ID       int64  `db:"id"`
	CourseID int64  `db:"course_id"`
	Title    string `db:"title"`
	Content  string `db:"content"`
}
