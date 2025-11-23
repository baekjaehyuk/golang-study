package model

import (
	"errors"
	"golang-course-registration/common/exception"
)

type Enrollment struct {
	ID        int `json:"id"`
	StudentID int `json:"student_id"`
	LectureID int `json:"lecture_id"`
}

func NewEnrollment(studentID, lectureID int) (*Enrollment, error) {
	if studentID < StudentIDMinLength || studentID > StudentIDMaxLength {
		return nil, errors.New(exception.ErrStudentIDInvalid)
	}

	if lectureID <= 0 {
		return nil, errors.New(exception.ErrEnrollmentLectureIDRequired)
	}

	return &Enrollment{
		StudentID: studentID,
		LectureID: lectureID,
	}, nil
}
