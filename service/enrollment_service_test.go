package service

import (
	"errors"
	"golang-course-registration/common/exception"
	"golang-course-registration/model"
	"strconv"
	"testing"
)

func TestEnrollmentService(t *testing.T) {
	t.Run("수강 신청", func(t *testing.T) {
		t.Run("성공", func(t *testing.T) {
			// given
			student, _ := model.NewStudent(1001)
			lecture, _ := model.NewLecture(2001, "데이터베이스", 30, 3, model.Monday, "09:00", "10:30")
			mockStudentRepo := &MockStudentRepositoryForService{students: []model.Student{*student}}
			mockLectureRepo := &MockLectureRepositoryForService{lectures: []model.Lecture{*lecture}}
			mockEnrollmentRepo := &MockEnrollmentRepositoryForService{enrollments: []model.Enrollment{}, lectures: []model.Lecture{*lecture}}
			service := NewEnrollmentService(mockEnrollmentRepo, mockLectureRepo, mockStudentRepo)

			// when
			response, err := service.Enroll(1001, 2001)

			// then
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if response.StudentID != 1001 || response.LectureID != 2001 {
				t.Errorf("expected (1001, 2001), got (%d, %d)", response.StudentID, response.LectureID)
			}
		})

		t.Run("예외 : 존재하지 않는 학생", func(t *testing.T) {
			// given
			lecture, _ := model.NewLecture(2001, "데이터베이스", 30, 3, model.Monday, "09:00", "10:30")
			mockStudentRepo := &MockStudentRepositoryForService{students: []model.Student{}}
			mockLectureRepo := &MockLectureRepositoryForService{lectures: []model.Lecture{*lecture}}
			mockEnrollmentRepo := &MockEnrollmentRepositoryForService{enrollments: []model.Enrollment{}, lectures: []model.Lecture{*lecture}}
			service := NewEnrollmentService(mockEnrollmentRepo, mockLectureRepo, mockStudentRepo)

			// when
			_, err := service.Enroll(1001, 2001)

			// then
			if err == nil || err.Error() != exception.ErrStudentNotFound {
				t.Errorf("expected error %s, got %v", exception.ErrStudentNotFound, err)
			}
		})

		t.Run("예외 : 존재하지 않는 강의", func(t *testing.T) {
			// given
			student, _ := model.NewStudent(1001)
			mockStudentRepo := &MockStudentRepositoryForService{students: []model.Student{*student}}
			mockLectureRepo := &MockLectureRepositoryForService{lectures: []model.Lecture{}}
			mockEnrollmentRepo := &MockEnrollmentRepositoryForService{enrollments: []model.Enrollment{}, lectures: []model.Lecture{}}
			service := NewEnrollmentService(mockEnrollmentRepo, mockLectureRepo, mockStudentRepo)

			// when
			_, err := service.Enroll(1001, 2001)

			// then
			if err == nil || err.Error() != exception.ErrLectureNotFound {
				t.Errorf("expected error %s, got %v", exception.ErrLectureNotFound, err)
			}
		})

		t.Run("예외 : 수강 정원 초과", func(t *testing.T) {
			// given
			student, _ := model.NewStudent(1001)
			lecture, _ := model.NewLecture(2001, "데이터베이스", 1, 3, model.Monday, "09:00", "10:30")
			lecture.CurrentEnrollment = 1
			mockStudentRepo := &MockStudentRepositoryForService{students: []model.Student{*student}}
			mockLectureRepo := &MockLectureRepositoryForService{lectures: []model.Lecture{*lecture}}
			mockEnrollmentRepo := &MockEnrollmentRepositoryForService{enrollments: []model.Enrollment{}, lectures: []model.Lecture{*lecture}}
			service := NewEnrollmentService(mockEnrollmentRepo, mockLectureRepo, mockStudentRepo)

			// when
			_, err := service.Enroll(1001, 2001)

			// then
			if err == nil || err.Error() != exception.ErrLectureCapacityExceeded {
				t.Errorf("expected error %s, got %v", exception.ErrLectureCapacityExceeded, err)
			}
		})

		t.Run("예외 : 강의 시간 충돌", func(t *testing.T) {
			// given
			student, _ := model.NewStudent(1001)
			existingLecture, _ := model.NewLecture(2001, "데이터베이스", 30, 3, model.Monday, "09:00", "10:30")
			newLecture, _ := model.NewLecture(2002, "운영체제", 30, 3, model.Monday, "10:00", "11:30")
			enrollment := model.Enrollment{StudentID: 1001, LectureID: 2001}
			mockStudentRepo := &MockStudentRepositoryForService{students: []model.Student{*student}}
			mockLectureRepo := &MockLectureRepositoryForService{lectures: []model.Lecture{*existingLecture, *newLecture}}
			mockEnrollmentRepo := &MockEnrollmentRepositoryForService{enrollments: []model.Enrollment{enrollment}, lectures: []model.Lecture{*existingLecture}}
			service := NewEnrollmentService(mockEnrollmentRepo, mockLectureRepo, mockStudentRepo)

			// when
			_, err := service.Enroll(1001, 2002)

			// then
			expectedError := exception.TimeConflictMessage(existingLecture.Name)
			if err == nil || err.Error() != expectedError {
				t.Errorf("expected error %s, got %v", expectedError, err)
			}
		})

		t.Run("예외 : 최대 수강 학점 초과", func(t *testing.T) {
			// given
			student, _ := model.NewStudent(1001)
			var lectures []model.Lecture
			var enrollments []model.Enrollment
			for i := 0; i < 6; i++ {
				lec, _ := model.NewLecture(2001+i, "강의"+strconv.Itoa(i), 30, 3, model.Monday, "09:00", "10:30")
				lectures = append(lectures, *lec)
				enrollments = append(enrollments, model.Enrollment{StudentID: 1001, LectureID: lec.ID})
			}
			newLecture, _ := model.NewLecture(3000, "추가 강의", 30, 3, model.Friday, "11:00", "12:30")
			lectures = append(lectures, *newLecture)

			mockStudentRepo := &MockStudentRepositoryForService{students: []model.Student{*student}}
			mockLectureRepo := &MockLectureRepositoryForService{lectures: lectures}
			mockEnrollmentRepo := &MockEnrollmentRepositoryForService{enrollments: enrollments, lectures: lectures[:6]}
			service := NewEnrollmentService(mockEnrollmentRepo, mockLectureRepo, mockStudentRepo)

			// when
			_, err := service.Enroll(1001, 3000)

			// then
			if err == nil || err.Error() != exception.ErrCreditLimitExceeded {
				t.Errorf("기대 : %s, 결과 : %v", exception.ErrCreditLimitExceeded, err)
			}
		})
	})

	t.Run("수강 신청 목록 조회", func(t *testing.T) {
		// given
		lecture1, _ := model.NewLecture(2001, "데이터베이스", 30, 3, model.Monday, "09:00", "10:30")
		lecture2, _ := model.NewLecture(2002, "운영체제", 30, 3, model.Tuesday, "09:00", "10:30")
		enrollments := []model.Enrollment{{StudentID: 1001, LectureID: 2001}, {StudentID: 1001, LectureID: 2002}}
		mockStudentRepo := &MockStudentRepositoryForService{}
		mockLectureRepo := &MockLectureRepositoryForService{lectures: []model.Lecture{*lecture1, *lecture2}}
		mockEnrollmentRepo := &MockEnrollmentRepositoryForService{enrollments: enrollments, lectures: []model.Lecture{*lecture1, *lecture2}}
		service := NewEnrollmentService(mockEnrollmentRepo, mockLectureRepo, mockStudentRepo)

		// when
		responses, _ := service.ListByStudent(1001)

		// then
		if len(responses) != 2 {
			t.Errorf("기대 : 2, 결과 : %d", len(responses))
		}
	})

	t.Run("수강 취소", func(t *testing.T) {
		t.Run("성공", func(t *testing.T) {
			// given
			student, _ := model.NewStudent(1001)
			lecture, _ := model.NewLecture(2001, "데이터베이스", 30, 3, model.Monday, "09:00", "10:30")
			lecture.CurrentEnrollment = 10
			enrollment := model.Enrollment{StudentID: 1001, LectureID: 2001}
			mockStudentRepo := &MockStudentRepositoryForService{students: []model.Student{*student}}
			mockLectureRepo := &MockLectureRepositoryForService{lectures: []model.Lecture{*lecture}}
			mockEnrollmentRepo := &MockEnrollmentRepositoryForService{enrollments: []model.Enrollment{enrollment}, lectures: []model.Lecture{*lecture}}
			service := NewEnrollmentService(mockEnrollmentRepo, mockLectureRepo, mockStudentRepo)

			// when
			_ = service.Cancel(1001, 2001)

			// then
			updatedLecture, _ := mockLectureRepo.FindByID(2001)
			if updatedLecture.CurrentEnrollment != 9 {
				t.Errorf("기대 : 9, 결과 : %d", updatedLecture.CurrentEnrollment)
			}
		})

		t.Run("예외 : 존재하지 않는 학생", func(t *testing.T) {
			// given
			lecture, _ := model.NewLecture(2001, "데이터베이스", 30, 3, model.Monday, "09:00", "10:30")
			mockStudentRepo := &MockStudentRepositoryForService{students: []model.Student{}}
			mockLectureRepo := &MockLectureRepositoryForService{lectures: []model.Lecture{*lecture}}
			mockEnrollmentRepo := &MockEnrollmentRepositoryForService{enrollments: []model.Enrollment{}, lectures: []model.Lecture{*lecture}}
			service := NewEnrollmentService(mockEnrollmentRepo, mockLectureRepo, mockStudentRepo)

			// when
			err := service.Cancel(1001, 2001)

			// then
			if err == nil || err.Error() != exception.ErrStudentNotFound {
				t.Errorf("기대 : %s, 결과 : %v", exception.ErrStudentNotFound, err)
			}
		})

		t.Run("예외 : 존재하지 않는 강의", func(t *testing.T) {
			// given
			student, _ := model.NewStudent(1001)
			mockStudentRepo := &MockStudentRepositoryForService{students: []model.Student{*student}}
			mockLectureRepo := &MockLectureRepositoryForService{lectures: []model.Lecture{}}
			mockEnrollmentRepo := &MockEnrollmentRepositoryForService{enrollments: []model.Enrollment{}, lectures: []model.Lecture{}}
			service := NewEnrollmentService(mockEnrollmentRepo, mockLectureRepo, mockStudentRepo)

			// when
			err := service.Cancel(1001, 2001)

			// then
			if err == nil || err.Error() != exception.ErrLectureNotFound {
				t.Errorf("기대 : %s, 결과 : %v", exception.ErrLectureNotFound, err)
			}
		})
	})
}

type MockEnrollmentRepositoryForService struct {
	enrollments []model.Enrollment
	lectures    []model.Lecture
	createError error
	deleteError error
}

func (m *MockEnrollmentRepositoryForService) Create(enrollment model.Enrollment) (model.Enrollment, error) {
	if m.createError != nil {
		return model.Enrollment{}, m.createError
	}
	enrollment.ID = len(m.enrollments) + 1
	m.enrollments = append(m.enrollments, enrollment)
	return enrollment, nil
}

func (m *MockEnrollmentRepositoryForService) FindByStudent(studentID int) ([]model.Enrollment, error) {
	var result []model.Enrollment
	for _, enrollment := range m.enrollments {
		if enrollment.StudentID == studentID {
			result = append(result, enrollment)
		}
	}
	return result, nil
}

func (m *MockEnrollmentRepositoryForService) FindLecturesByStudent(studentID int) ([]model.Lecture, error) {
	var result []model.Lecture
	for _, enrollment := range m.enrollments {
		if enrollment.StudentID == studentID {
			for _, lecture := range m.lectures {
				if lecture.ID == enrollment.LectureID {
					result = append(result, lecture)
				}
			}
		}
	}
	return result, nil
}

func (m *MockEnrollmentRepositoryForService) CountByLectureID(lectureID int) (int, error) {
	count := 0
	for _, enrollment := range m.enrollments {
		if enrollment.LectureID == lectureID {
			count++
		}
	}
	return count, nil
}

func (m *MockEnrollmentRepositoryForService) DeleteByStudentAndLecture(studentID, lectureID int) error {
	if m.deleteError != nil {
		return m.deleteError
	}
	for i, enrollment := range m.enrollments {
		if enrollment.StudentID == studentID && enrollment.LectureID == lectureID {
			m.enrollments = append(m.enrollments[:i], m.enrollments[i+1:]...)
			return nil
		}
	}
	return errors.New("enrollment not found")
}

type MockLectureRepositoryForService struct {
	lectures      []model.Lecture
	findByIDError error
	updateError   error
}

func (m *MockLectureRepositoryForService) FindAll() ([]model.Lecture, error) {
	return m.lectures, nil
}

func (m *MockLectureRepositoryForService) FindByID(id int) (model.Lecture, error) {
	if m.findByIDError != nil {
		return model.Lecture{}, m.findByIDError
	}
	for _, lecture := range m.lectures {
		if lecture.ID == id {
			return lecture, nil
		}
	}
	return model.Lecture{}, errors.New(exception.ErrLectureNotFound)
}

func (m *MockLectureRepositoryForService) FindByName(name string) (model.Lecture, error) {
	for _, lecture := range m.lectures {
		if lecture.Name == name {
			return lecture, nil
		}
	}
	return model.Lecture{}, errors.New(exception.ErrLectureNotFound)
}

func (m *MockLectureRepositoryForService) Create(lecture model.Lecture) (model.Lecture, error) {
	m.lectures = append(m.lectures, lecture)
	return lecture, nil
}

func (m *MockLectureRepositoryForService) Delete(id int) error {
	for i, lecture := range m.lectures {
		if lecture.ID == id {
			m.lectures = append(m.lectures[:i], m.lectures[i+1:]...)
			return nil
		}
	}
	return errors.New(exception.ErrLectureNotFound)
}

func (m *MockLectureRepositoryForService) UpdateCurrentEnrollment(lectureID int, currentEnrollment int) error {
	if m.updateError != nil {
		return m.updateError
	}
	for i, lecture := range m.lectures {
		if lecture.ID == lectureID {
			m.lectures[i].CurrentEnrollment = currentEnrollment
			return nil
		}
	}
	return errors.New(exception.ErrLectureNotFound)
}

type MockStudentRepositoryForService struct {
	students      []model.Student
	findByIDError error
}

func (m *MockStudentRepositoryForService) Create(student model.Student) (model.Student, error) {
	m.students = append(m.students, student)
	return student, nil
}

func (m *MockStudentRepositoryForService) FindByID(id int) (model.Student, error) {
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
