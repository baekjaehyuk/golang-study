package model

import (
	"golang-course-registration/common/exception"
	"testing"
)

func TestNewEnrollment(t *testing.T) {
	t.Run("성공", func(t *testing.T) {
		// given
		studentID := 1234
		lectureID := 5678

		// when
		enrollment, _ := NewEnrollment(studentID, lectureID)

		// then
		if enrollment.StudentID != studentID {
			t.Errorf("기대 : %d, 결과 : %d", studentID, enrollment.StudentID)
		}
		if enrollment.LectureID != lectureID {
			t.Errorf("기대 : %d, 결과ㅜ: %d", lectureID, enrollment.LectureID)
		}
	})

	t.Run("예외 : 유효하지 않은 학생 ID", func(t *testing.T) {
		// given
		invalidStudentID := 999
		lectureID := 5678

		// when
		_, err := NewEnrollment(invalidStudentID, lectureID)

		// then
		if err.Error() != exception.ErrStudentIDInvalid {
			t.Errorf("기대 : %s, 결과 : %v", exception.ErrStudentIDInvalid, err)
		}
	})

	t.Run("예외 : 유효하지 않은 강의 ID", func(t *testing.T) {
		// given
		studentID := 1234
		invalidLectureID := 0

		// when
		_, err := NewEnrollment(studentID, invalidLectureID)

		// then
		if err.Error() != exception.ErrEnrollmentLectureIDRequired {
			t.Errorf("기대 : %s, 결과 : %v", exception.ErrEnrollmentLectureIDRequired, err)
		}
	})
}
