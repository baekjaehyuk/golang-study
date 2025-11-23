package repository

import (
	"errors"
	"golang-course-registration/common/exception"
	"golang-course-registration/model"
	"strconv"

	"github.com/supabase-community/postgrest-go"
	"github.com/supabase-community/supabase-go"
)

type LectureRepository interface {
	FindAll() ([]model.Lecture, error)
	FindByID(id int) (model.Lecture, error)
	FindByName(name string) (model.Lecture, error)
	Create(lecture model.Lecture) (model.Lecture, error)
	Delete(id int) error
	UpdateCurrentEnrollment(lectureID int, currentEnrollment int) error
}

type lectureRepository struct {
	client *supabase.Client
}

func NewLectureRepository(client *supabase.Client) LectureRepository {
	return &lectureRepository{client: client}
}

func (r *lectureRepository) FindAll() ([]model.Lecture, error) {
	var result []model.Lecture
	_, err := r.client.From("lectures").
		Select("*", "", false).
		Order("id", &postgrest.OrderOpts{Ascending: true}).
		ExecuteTo(&result)
	return result, err
}

func (r *lectureRepository) FindByID(id int) (model.Lecture, error) {
	var result []model.Lecture
	_, err := r.client.From("lectures").
		Select("*", "", false).
		Eq("id", strconv.Itoa(id)).
		Limit(1, "").
		ExecuteTo(&result)
	if err != nil {
		return model.Lecture{}, err
	}

	if len(result) == 0 {
		return model.Lecture{}, errors.New(exception.ErrLectureNotFound)
	}

	return result[0], nil
}

func (r *lectureRepository) FindByName(name string) (model.Lecture, error) {
	var result []model.Lecture
	_, err := r.client.From("lectures").
		Select("*", "", false).
		Eq("name", name).
		Limit(1, "").
		ExecuteTo(&result)

	if err != nil {
		return model.Lecture{}, err
	}

	if len(result) == 0 {
		return model.Lecture{}, errors.New(exception.ErrLectureNotFound)
	}

	return result[0], nil
}

func (r *lectureRepository) Create(lecture model.Lecture) (model.Lecture, error) {
	var result []model.Lecture
	_, err := r.client.From("lectures").
		Insert(lecture, false, "", "representation", "").
		ExecuteTo(&result)
	if err != nil {
		return model.Lecture{}, err
	}

	if len(result) == 0 {
		return model.Lecture{}, errors.New("강좌 생성 결과가 비어 있습니다")
	}

	return result[0], nil
}

func (r *lectureRepository) Delete(id int) error {
	_, _, err := r.client.From("lectures").
		Delete("", "").
		Eq("id", strconv.Itoa(id)).
		Execute()
	return err
}

func (r *lectureRepository) UpdateCurrentEnrollment(lectureID int, currentEnrollment int) error {
	updateData := map[string]interface{}{
		"current_enrollment": currentEnrollment,
	}

	_, _, err := r.client.From("lectures").
		Update(updateData, "", "").
		Eq("id", strconv.Itoa(lectureID)).
		Execute()
	return err
}
