package models

import "time"

type Review struct {
	ID        int64     `db:"id"`
	StudentID int64     `db:"student_id"`
	CourseID  int64     `db:"course_id"`
	Comment   string    `db:"comment"`
	Rating    int32     `db:"rating"`
	CreatedAt time.Time `db:"created_at"`
}
