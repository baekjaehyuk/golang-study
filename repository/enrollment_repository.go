package repository

import (
	"golang-course-registration/model"
	"strconv"

	"github.com/supabase-community/postgrest-go"
	"github.com/supabase-community/supabase-go"
)

type EnrollmentRepository interface {
	Create(enrollment model.Enrollment) (model.Enrollment, error)
	FindByStudent(studentID int) ([]model.Enrollment, error)
	FindLecturesByStudent(studentID int) ([]model.Lecture, error)
	CountByLectureID(lectureID int) (int, error)
	DeleteByStudentAndLecture(studentID, lectureID int) error
}

type enrollmentRepository struct {
	client *supabase.Client
}

type enrollmentRecord struct {
	ID        int `json:"id"`
	StudentID int `json:"student_id"`
	LectureID int `json:"lecture_id"`
}

func (er enrollmentRecord) toModel() model.Enrollment {
	return model.Enrollment{
		ID:        er.ID,
		StudentID: er.StudentID,
		LectureID: er.LectureID,
	}
}

func NewEnrollmentRepository(client *supabase.Client) EnrollmentRepository {
	return &enrollmentRepository{client: client}
}

func (r *enrollmentRepository) Create(enrollment model.Enrollment) (model.Enrollment, error) {
	payload := map[string]interface{}{
		"student_id": enrollment.StudentID,
		"lecture_id": enrollment.LectureID,
	}

	var inserted []enrollmentRecord
	_, err := r.client.From("enrollments").
		Insert(payload, false, "", "representation", "").
		ExecuteTo(&inserted)
	if err != nil {
		return model.Enrollment{}, err
	}

	return inserted[0].toModel(), nil
}

func (r *enrollmentRepository) FindByStudent(studentID int) ([]model.Enrollment, error) {
	var records []enrollmentRecord
	_, err := r.client.From("enrollments").
		Select("*", "", false).
		Eq("student_id", strconv.Itoa(studentID)).
		Order("id", &postgrest.OrderOpts{Ascending: false}).
		ExecuteTo(&records)
	if err != nil {
		return nil, err
	}

	list := make([]model.Enrollment, 0, len(records))
	for _, record := range records {
		list = append(list, record.toModel())
	}
	return list, nil
}

func (r *enrollmentRepository) FindLecturesByStudent(studentID int) ([]model.Lecture, error) {
	var lectures []model.Lecture
	_, err := r.client.From("lectures").
		Select("*, enrollments!inner(student_id)", "", false).
		Eq("enrollments.student_id", strconv.Itoa(studentID)).
		Order("id", &postgrest.OrderOpts{Ascending: true}).
		ExecuteTo(&lectures)
	return lectures, err
}

func (r *enrollmentRepository) CountByLectureID(lectureID int) (int, error) {
	var records []enrollmentRecord
	_, err := r.client.From("enrollments").
		Select("*", "", false).
		Eq("lecture_id", strconv.Itoa(lectureID)).
		ExecuteTo(&records)
	if err != nil {
		return 0, err
	}
	return len(records), nil
}

func (r *enrollmentRepository) DeleteByStudentAndLecture(studentID, lectureID int) error {
	_, _, err := r.client.From("enrollments").
		Delete("", "").
		Eq("student_id", strconv.Itoa(studentID)).
		Eq("lecture_id", strconv.Itoa(lectureID)).
		Execute()
	return err
}
