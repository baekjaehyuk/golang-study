package service

import (
	"errors"
	"golang-course-registration/common/exception"
	"golang-course-registration/controller/dto"
	"golang-course-registration/model"
	"testing"
)

func TestLectureService(t *testing.T) {
	t.Run("강좌 개설", func(t *testing.T) {
		t.Run("성공", func(t *testing.T) {
			// given
			mockRepo := &MockLectureRepository{lectures: []model.Lecture{}}
			service := NewLectureService(mockRepo)
			req := dto.CreateLectureRequest{
				ID: 1001, Name: "데이터베이스", Capacity: 30, Credit: 3,
				Day: model.Monday, StartTime: "09:00", EndTime: "10:30",
			}

			// when
			response, _ := service.Create(req)

			// then
			if response.ID != 1001 || response.Name != "데이터베이스" {
				t.Errorf("기대 : (1001, 데이터베이스), 결과 : (%d, %s)", response.ID, response.Name)
			}
		})

		t.Run("예외 : 중복된 강좌명", func(t *testing.T) {
			// given
			existingLecture, _ := model.NewLecture(1000, "데이터베이스", 30, 3, model.Monday, "09:00", "10:30")
			mockRepo := &MockLectureRepository{lectures: []model.Lecture{*existingLecture}}
			service := NewLectureService(mockRepo)
			req := dto.CreateLectureRequest{
				ID: 1001, Name: "데이터베이스", Capacity: 30, Credit: 3,
				Day: model.Monday, StartTime: "09:00", EndTime: "10:30",
			}

			// when
			_, err := service.Create(req)

			// then
			if err == nil || err.Error() != exception.ErrLectureNameDuplicate {
				t.Errorf("기대 : %s, 결과 : %v", exception.ErrLectureNameDuplicate, err)
			}
		})

		t.Run("예외 : 중복된 강좌 번호", func(t *testing.T) {
			// given
			existingLecture, _ := model.NewLecture(1001, "운영체제", 30, 3, model.Monday, "09:00", "10:30")
			mockRepo := &MockLectureRepository{lectures: []model.Lecture{*existingLecture}}
			service := NewLectureService(mockRepo)
			req := dto.CreateLectureRequest{
				ID: 1001, Name: "데이터베이스", Capacity: 30, Credit: 3,
				Day: model.Monday, StartTime: "09:00", EndTime: "10:30",
			}

			// when
			_, err := service.Create(req)

			// then
			if err == nil || err.Error() != exception.ErrLectureIDDuplicate {
				t.Errorf("기대 : %s, 결과 : %v", exception.ErrLectureIDDuplicate, err)
			}
		})
	})

	t.Run("강좌 조회", func(t *testing.T) {
		t.Run("조회 성공", func(t *testing.T) {
			// given
			lecture, _ := model.NewLecture(1001, "데이터베이스", 30, 3, model.Monday, "09:00", "10:30")
			mockRepo := &MockLectureRepository{lectures: []model.Lecture{*lecture}}
			service := NewLectureService(mockRepo)

			// when
			response, _ := service.FindByID(1001)

			// then
			if response.ID != 1001 {
				t.Errorf("기대 : 1001, 결과: %d", response.ID)
			}
		})

		t.Run("예외: 존재하지 않는 강좌", func(t *testing.T) {
			// given
			mockRepo := &MockLectureRepository{lectures: []model.Lecture{}}
			service := NewLectureService(mockRepo)

			// when
			_, err := service.FindByID(9999)

			// then
			if err == nil || err.Error() != exception.ErrLectureNotFound {
				t.Errorf("기대 : %s, 결과 : %v", exception.ErrLectureNotFound, err)
			}
		})

		t.Run("전체 목록 조회", func(t *testing.T) {
			// given
			lecture1, _ := model.NewLecture(1001, "데이터베이스", 30, 3, model.Monday, "09:00", "10:30")
			lecture2, _ := model.NewLecture(1002, "운영체제", 25, 3, model.Tuesday, "11:00", "12:30")
			mockRepo := &MockLectureRepository{lectures: []model.Lecture{*lecture1, *lecture2}}
			service := NewLectureService(mockRepo)

			// when
			responses, _ := service.List()

			// then
			if len(responses) != 2 {
				t.Errorf("기대 : 2, 결과 : %d", len(responses))
			}
		})
	})

	t.Run("강좌 삭제", func(t *testing.T) {
		t.Run("성공", func(t *testing.T) {
			// given
			lecture, _ := model.NewLecture(1001, "데이터베이스", 30, 3, model.Monday, "09:00", "10:30")
			mockRepo := &MockLectureRepository{lectures: []model.Lecture{*lecture}}
			service := NewLectureService(mockRepo)

			// when
			_ = service.Delete(1001)

			// then
			_, findErr := service.FindByID(1001)
			if findErr == nil || findErr.Error() != exception.ErrLectureNotFound {
				t.Error("강의가 삭제되지 않았습니다.")
			}
		})

		t.Run("예외 : 존재하지 않는 강좌", func(t *testing.T) {
			// given
			mockRepo := &MockLectureRepository{lectures: []model.Lecture{}}
			service := NewLectureService(mockRepo)

			// when
			err := service.Delete(9999)

			// then
			if err == nil || err.Error() != exception.ErrLectureNotFound {
				t.Errorf("기대 : %s, 결과 : %v", exception.ErrLectureNotFound, err)
			}
		})
	})
}

type MockLectureRepository struct {
	lectures        []model.Lecture
	findByIDError   error
	findByNameError error
	createError     error
	deleteError     error
	updateError     error
}

func (m *MockLectureRepository) FindAll() ([]model.Lecture, error) {
	return m.lectures, nil
}

func (m *MockLectureRepository) FindByID(id int) (model.Lecture, error) {
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

func (m *MockLectureRepository) FindByName(name string) (model.Lecture, error) {
	if m.findByNameError != nil {
		return model.Lecture{}, m.findByNameError
	}
	for _, lecture := range m.lectures {
		if lecture.Name == name {
			return lecture, nil
		}
	}
	return model.Lecture{}, errors.New(exception.ErrLectureNotFound)
}

func (m *MockLectureRepository) Create(lecture model.Lecture) (model.Lecture, error) {
	if m.createError != nil {
		return model.Lecture{}, m.createError
	}
	m.lectures = append(m.lectures, lecture)
	return lecture, nil
}

func (m *MockLectureRepository) Delete(id int) error {
	if m.deleteError != nil {
		return m.deleteError
	}
	for i, lecture := range m.lectures {
		if lecture.ID == id {
			m.lectures = append(m.lectures[:i], m.lectures[i+1:]...)
			return nil
		}
	}
	return errors.New(exception.ErrLectureNotFound)
}

func (m *MockLectureRepository) UpdateCurrentEnrollment(lectureID int, currentEnrollment int) error {
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

type MockEnrollmentRepository struct {
	enrollments []model.Enrollment
	lectures    []model.Lecture
}

func (m *MockEnrollmentRepository) Create(enrollment model.Enrollment) (model.Enrollment, error) {
	enrollment.ID = len(m.enrollments) + 1
	m.enrollments = append(m.enrollments, enrollment)
	return enrollment, nil
}

func (m *MockEnrollmentRepository) FindByStudent(studentID int) ([]model.Enrollment, error) {
	var result []model.Enrollment
	for _, enrollment := range m.enrollments {
		if enrollment.StudentID == studentID {
			result = append(result, enrollment)
		}
	}
	return result, nil
}

func (m *MockEnrollmentRepository) FindLecturesByStudent(studentID int) ([]model.Lecture, error) {
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

func (m *MockEnrollmentRepository) CountByLectureID(lectureID int) (int, error) {
	count := 0
	for _, enrollment := range m.enrollments {
		if enrollment.LectureID == lectureID {
			count++
		}
	}
	return count, nil
}

func (m *MockEnrollmentRepository) DeleteByStudentAndLecture(studentID, lectureID int) error {
	for i, enrollment := range m.enrollments {
		if enrollment.StudentID == studentID && enrollment.LectureID == lectureID {
			m.enrollments = append(m.enrollments[:i], m.enrollments[i+1:]...)
			return nil
		}
	}
	return errors.New("enrollment not found")
}
