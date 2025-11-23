package repository

import (
	"errors"
	"golang-course-registration/model"
	"strconv"

	"github.com/supabase-community/supabase-go"
)

type StudentRepository interface {
	Create(student model.Student) (model.Student, error)
	FindByID(id int) (model.Student, error)
}

type studentRepository struct {
	client *supabase.Client
}

func NewStudentRepository(client *supabase.Client) StudentRepository {
	return &studentRepository{client: client}
}

func (r *studentRepository) Create(student model.Student) (model.Student, error) {
	_, _, err := r.client.From("students").
		Insert(student, false, "", "minimal", "").
		Execute()
	if err != nil {
		return model.Student{}, err
	}
	return student, nil
}

func (r *studentRepository) FindByID(id int) (model.Student, error) {
	var list []model.Student
	_, err := r.client.From("students").
		Select("*", "", false).
		Eq("id", strconv.Itoa(id)).
		Limit(1, "").
		ExecuteTo(&list)
	if err != nil {
		return model.Student{}, err
	}
	if len(list) == 0 {
		return model.Student{}, errors.New("student not found")
	}
	return list[0], nil
}
