package model

import (
	"errors"
	"golang-course-registration/common/exception"
)

const (
	StudentIDMinLength = 1000
	StudentIDMaxLength = 9999
)

type Student struct {
	ID int `json:"id"`
}

func NewStudent(id int) (*Student, error) {
	if id < StudentIDMinLength || id > StudentIDMaxLength {
		return nil, errors.New(exception.ErrStudentIDInvalid)
	}

	return &Student{ID: id}, nil
}
