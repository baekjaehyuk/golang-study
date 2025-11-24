package model

import (
	"golang-course-registration/common/constants"
	"golang-course-registration/common/exception"
	"testing"
)

func TestNewStudent(t *testing.T) {
	// given
	t.Run("성공 : 유효한 학번으로 학생 생성", func(t *testing.T) {
		// when
		student, _ := NewStudent(1234)
		// then
		if student == nil {
			t.Error("학생이 생성되지 않았습니다.")
		}
		if student.ID != 1234 {
			t.Errorf("기대값: 1234, 실제값: %d", student.ID)
		}
	})

	// given
	t.Run("실패 : 유효하지 않은 학번으로 학생 생성", func(t *testing.T) {
		tests := []struct {
			name string
			id   int
		}{
			{"학번이 1000보다 작을 경우 예외 발생", constants.StudentIdMin - 1},
			{"학번이 9999보다 클 경우 예외 발생", constants.StudentIdMax + 1},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				// when
				_, err := NewStudent(tt.id)
				// then
				if err == nil {
					t.Error("오류가 발생해야 합니다.")
				}
				if err.Error() != exception.ErrStudentIDInvalid {
					t.Errorf("기대 오류: %s, 실제 오류: %v", exception.ErrStudentIDInvalid, err)
				}
			})
		}
	})
}
