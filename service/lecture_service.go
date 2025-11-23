package service

import (
	"errors"
	"golang-course-registration/common/exception"
	"golang-course-registration/controller/dto"
	"golang-course-registration/model"
	"golang-course-registration/repository"
)

type LectureService interface {
	Create(req dto.CreateLectureRequest) (dto.LectureResponse, error)
	FindByID(id int) (dto.LectureResponse, error)
	List() ([]dto.LectureResponse, error)
	Delete(id int) error
}

type lectureService struct {
	lectureRepo    repository.LectureRepository
	enrollmentRepo repository.EnrollmentRepository
}

func NewLectureService(lectureRepo repository.LectureRepository) LectureService {
	return &lectureService{lectureRepo: lectureRepo}
}

func NewLectureServiceWithEnrollment(lectureRepo repository.LectureRepository, enrollmentRepo repository.EnrollmentRepository) LectureService {
	return &lectureService{
		lectureRepo:    lectureRepo,
		enrollmentRepo: enrollmentRepo,
	}
}

func (s *lectureService) Create(req dto.CreateLectureRequest) (dto.LectureResponse, error) {
	lecture, err := model.NewLecture(
		req.ID,
		req.Name,
		req.Capacity,
		req.Credit,
		req.Day,
		req.StartTime,
		req.EndTime,
	)

	if err != nil {
		return dto.LectureResponse{}, err
	}

	_, errExistName := s.lectureRepo.FindByName(lecture.Name)
	if errExistName == nil {
		return dto.LectureResponse{}, errors.New(exception.ErrLectureNameDuplicate)
	}

	_, errExistID := s.lectureRepo.FindByID(lecture.ID)
	if errExistID == nil {
		return dto.LectureResponse{}, errors.New(exception.ErrLectureIDDuplicate)
	}

	createdLecture, err := s.lectureRepo.Create(*lecture)
	if err != nil {
		return dto.LectureResponse{}, err
	}

	return dto.NewLectureResponse(createdLecture), nil
}

func (s *lectureService) FindByID(id int) (dto.LectureResponse, error) {
	lecture, err := s.lectureRepo.FindByID(id)
	if err != nil {
		return dto.LectureResponse{}, err
	}
	return dto.NewLectureResponse(lecture), nil
}

func (s *lectureService) List() ([]dto.LectureResponse, error) {
	lectures, err := s.lectureRepo.FindAll()
	if err != nil {
		return nil, err
	}

	responses := make([]dto.LectureResponse, 0, len(lectures))
	for _, lecture := range lectures {
		response := dto.NewLectureResponse(lecture)
		responses = append(responses, response)
	}

	return responses, nil
}

func (s *lectureService) Delete(id int) error {
	_, err := s.lectureRepo.FindByID(id)
	if err != nil {
		return errors.New(exception.ErrLectureNotFound)
	}

	return s.lectureRepo.Delete(id)
}
