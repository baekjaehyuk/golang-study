package service

import (
	"errors"
	"golang-course-registration/common/exception"
	"golang-course-registration/model"
	"testing"
)

func TestStudentService(t *testing.T) {
	t.Run("학생 등록", func(t *testing.T) {
		t.Run("성공", func(t *testing.T) {
			// given
			mockRepo := &MockStudentRepository{students: []model.Student{}}
			service := NewStudentService(mockRepo)

			// when
			response, _ := service.Register(1001)

			// then
			if response.ID != 1001 {
				t.Errorf("기대 : 1001, 결과 :  %d", response.ID)
			}
		})

		t.Run("예외 : 유효하지 않은 학번", func(t *testing.T) {
			// given
			mockRepo := &MockStudentRepository{students: []model.Student{}}
			service := NewStudentService(mockRepo)
			testCases := []struct {
				name string
				id   int
			}{
				{"학번이 너무 작은 경우", 999},
				{"학번이 너무 큰 경우", 10000},
			}

			for _, tc := range testCases {
				t.Run(tc.name, func(t *testing.T) {
					// when
					_, err := service.Register(tc.id)

					// then
					if err.Error() != exception.ErrStudentIDInvalid {
						t.Errorf("기대 : %s, 결과 : %s", exception.ErrStudentIDInvalid, err.Error())
					}
				})
			}
		})
	})
}

type MockStudentRepository struct {
	students      []model.Student
	findByIDError error
	createError   error
}

func (m *MockStudentRepository) Create(student model.Student) (model.Student, error) {
	if m.createError != nil {
		return model.Student{}, m.createError
	}
	for _, s := range m.students {
		if s.ID == student.ID {
			return student, nil
		}
	}
	m.students = append(m.students, student)
	return student, nil
}

func (m *MockStudentRepository) FindByID(id int) (model.Student, error) {
	if m.findByIDError != nil {
		return model.Student{}, m.findByIDError
	}
	for _, student := range m.students {
		if student.ID == id {
			return student, nil
		}
	}
	return model.Student{}, errors.New(exception.ErrStudentNotFound)
}
