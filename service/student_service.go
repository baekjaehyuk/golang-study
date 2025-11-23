package service

import (
	"golang-course-registration/controller/dto"
	"golang-course-registration/model"
	"golang-course-registration/repository"
)

type StudentService interface {
	Register(studentId int) (dto.StudentResponse, error)
}

type studentService struct {
	repo repository.StudentRepository
}

func NewStudentService(repo repository.StudentRepository) StudentService {
	return &studentService{repo: repo}
}

func (s *studentService) Register(id int) (dto.StudentResponse, error) {
	student, err := model.NewStudent(id)
	if err != nil {
		return dto.StudentResponse{}, err
	}

	savedStudent, err := s.repo.Create(*student)
	if err != nil {
		return dto.StudentResponse{}, err
	}

	return dto.NewStudentResponse(savedStudent), nil
}
