package service

import (
	"errors"
	"golang-course-registration/common/exception"
	"golang-course-registration/controller/dto"
	"golang-course-registration/model"
	"golang-course-registration/repository"
	"sync"
	"time"
)

type EnrollmentService interface {
	Enroll(studentID, lectureID int) (dto.EnrollmentResponse, error)
	Cancel(studentID, lectureID int) error
	ListByStudent(studentID int) ([]dto.LectureResponse, error)
}

type enrollmentService struct {
	enrollmentRepo repository.EnrollmentRepository
	lectureRepo    repository.LectureRepository
	studentRepo    repository.StudentRepository
	lectureLocks   map[int]*sync.Mutex
	locksMutex     sync.Mutex
}

func NewEnrollmentService(
	enrollmentRepo repository.EnrollmentRepository,
	lectureRepo repository.LectureRepository,
	studentRepo repository.StudentRepository,
) EnrollmentService {
	return &enrollmentService{
		enrollmentRepo: enrollmentRepo,
		lectureRepo:    lectureRepo,
		studentRepo:    studentRepo,
		lectureLocks:   make(map[int]*sync.Mutex),
	}
}

// Enroll 수강신청
func (s *enrollmentService) Enroll(studentID, lectureID int) (dto.EnrollmentResponse, error) {
	lectureLock := s.getLectureLock(lectureID)
	lectureLock.Lock()
	defer lectureLock.Unlock()

	lecture, err := s.validateEnrollment(studentID, lectureID)
	if err != nil {
		return dto.EnrollmentResponse{}, err
	}

	if err := s.checkTimeConflict(studentID, lecture); err != nil {
		return dto.EnrollmentResponse{}, err
	}

	if err := s.checkCreditLimit(studentID, lecture); err != nil {
		return dto.EnrollmentResponse{}, err
	}

	return s.createEnrollment(studentID, lectureID)
}

// ListByStudent 학생 수강신청 내역 조회
func (s *enrollmentService) ListByStudent(studentID int) ([]dto.LectureResponse, error) {
	lectures, err := s.enrollmentRepo.FindLecturesByStudent(studentID)
	if err != nil {
		return nil, err
	}

	lectureList := make([]dto.LectureResponse, 0, len(lectures))
	for _, lecture := range lectures {
		response := dto.NewLectureResponse(lecture)
		lectureList = append(lectureList, response)
	}

	return lectureList, nil
}

// validateEnrollment 학생 및 강좌 존재 여부, 정원 체크
func (s *enrollmentService) validateEnrollment(studentID, lectureID int) (model.Lecture, error) {
	if _, err := s.studentRepo.FindByID(studentID); err != nil {
		return model.Lecture{}, errors.New(exception.ErrStudentNotFound)
	}

	lecture, err := s.lectureRepo.FindByID(lectureID)
	if err != nil {
		return model.Lecture{}, errors.New(exception.ErrLectureNotFound)
	}

	if lecture.IsFull() {
		return model.Lecture{}, errors.New(exception.ErrLectureCapacityExceeded)
	}

	return lecture, nil
}

// checkTimeConflict 기존 수강신청과 시간 충돌 체크
func (s *enrollmentService) checkTimeConflict(studentID int, newLecture model.Lecture) error {
	existingLectures, err := s.enrollmentRepo.FindLecturesByStudent(studentID)
	if err != nil {
		return err
	}

	for _, existingLecture := range existingLectures {
		if hasTimeConflict(newLecture, existingLecture) {
			return errors.New(exception.TimeConflictMessage(existingLecture.Name))
		}
	}

	return nil
}

// checkCreditLimit 총 학점이 18학점을 초과하지 않는지 체크
func (s *enrollmentService) checkCreditLimit(studentID int, newLecture model.Lecture) error {
	existingLectures, err := s.enrollmentRepo.FindLecturesByStudent(studentID)
	if err != nil {
		return err
	}

	totalCredit := newLecture.Credit
	for _, existingLecture := range existingLectures {
		totalCredit += existingLecture.Credit
	}

	if totalCredit > 18 {
		return errors.New(exception.ErrCreditLimitExceeded)
	}

	return nil
}

// createEnrollment 수강신청 생성 및 현재 수강 인원 증가
func (s *enrollmentService) createEnrollment(studentID, lectureID int) (dto.EnrollmentResponse, error) {
	enrollment, err := model.NewEnrollment(studentID, lectureID)
	if err != nil {
		return dto.EnrollmentResponse{}, err
	}

	createdEnrollment, err := s.enrollmentRepo.Create(*enrollment)
	if err != nil {
		return dto.EnrollmentResponse{}, err
	}

	// 강좌 현재 수강 인원 업데이트
	lecture, err := s.lectureRepo.FindByID(lectureID)
	if err != nil {
		return dto.EnrollmentResponse{}, err
	}
	lecture.IncrementEnrollment()
	if err := s.lectureRepo.UpdateCurrentEnrollment(lectureID, lecture.CurrentEnrollment); err != nil {
		return dto.EnrollmentResponse{}, err
	}

	return dto.NewEnrollmentResponse(createdEnrollment), nil
}

// Cancel 수강신청 취소
func (s *enrollmentService) Cancel(studentID, lectureID int) error {
	lectureLock := s.getLectureLock(lectureID)
	lectureLock.Lock()
	defer lectureLock.Unlock()

	if _, err := s.studentRepo.FindByID(studentID); err != nil {
		return errors.New(exception.ErrStudentNotFound)
	}

	lecture, err := s.lectureRepo.FindByID(lectureID)
	if err != nil {
		return errors.New(exception.ErrLectureNotFound)
	}

	if err := s.enrollmentRepo.DeleteByStudentAndLecture(studentID, lectureID); err != nil {
		return err
	}

	// 강좌 현재 수강 인원 업데이트
	lecture.DecrementEnrollment()
	if err := s.lectureRepo.UpdateCurrentEnrollment(lectureID, lecture.CurrentEnrollment); err != nil {
		return err
	}

	return nil
}

func hasTimeConflict(newLec, existLec model.Lecture) bool {
	if newLec.Day != existLec.Day {
		return false
	}
	newLecStart, _ := time.Parse("15:04", newLec.StartTime)
	newLecEnd, _ := time.Parse("15:04", newLec.EndTime)
	existLecStart, _ := time.Parse("15:04", existLec.StartTime)
	existLecEnd, _ := time.Parse("15:04", existLec.EndTime)
	return newLecStart.Before(existLecEnd) && newLecEnd.After(existLecStart)
}

func (s *enrollmentService) getLectureLock(lectureID int) *sync.Mutex {
	s.locksMutex.Lock()
	defer s.locksMutex.Unlock()

	lock, exists := s.lectureLocks[lectureID]
	if !exists {
		lock = &sync.Mutex{}
		s.lectureLocks[lectureID] = lock
	}
	return lock
}
