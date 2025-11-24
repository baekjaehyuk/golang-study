package model

import (
	"golang-course-registration/common/exception"
	"testing"
)

func TestNewLecture(t *testing.T) {
	// given
	t.Run("강좌 생성", func(t *testing.T) {
		// when
		lecture, _ := NewLecture(1001, "자료구조", 20, 3, Monday, "10:00", "13:00")
		if lecture == nil {
			t.Error("강의가 생성되지 않았습니다.")
		}
	})

	// given
	t.Run("예외 : 긴 한글 강좌명으로 생성", func(t *testing.T) {
		longName := "데이터베이스설계와구조및기초프로그래밍의이해"
		// when
		lecture, err := NewLecture(1002, longName, 20, 3, Tuesday, "14:00", "17:00")
		// then
		if err == nil {
			t.Errorf("기대 : %s, 결과 : %s", longName, lecture.Name)
		}
	})

	// given
	t.Run("예외 : 강좌 번호가 작은 경우", func(t *testing.T) {
		// when
		_, err := NewLecture(999, "운영체제", 25, 3, Wednesday, "09:00", "12:00")
		// then
		if err == nil || err.Error() != exception.ErrLectureIDInvalid {
			t.Errorf("기대 : %s, 결과 : %v", exception.ErrLectureIDInvalid, err)
		}
	})

	// given
	t.Run("예외 : 강좌명이 짧은 경우", func(t *testing.T) {
		// when
		_, err := NewLecture(1003, "A", 20, 3, Thursday, "13:00", "15:00")
		// then
		if err == nil || err.Error() != exception.ErrLectureNameRequired {
			t.Errorf("기대 : %s, 결과 : %v", exception.ErrLectureNameRequired, err)
		}
	})

	// given
	t.Run("예외 : 수강 인원이 많은 경우", func(t *testing.T) {
		// when
		_, err := NewLecture(1004, "네트워크", 50, 3, Friday, "10:00", "12:00")
		// then
		if err == nil || err.Error() != exception.ErrLectureCapacityInvalid {
			t.Errorf("기대 : %s, 결과 : %v", exception.ErrLectureCapacityInvalid, err)
		}
	})

	// given
	t.Run("예외 : 학점이 유효하지 않은 경우", func(t *testing.T) {
		// when
		_, err := NewLecture(1005, "컴퓨터 구조", 30, 7, Monday, "15:00", "18:00")
		// then
		if err == nil || err.Error() != exception.ErrLectureCreditInvalid {
			t.Errorf("기대 오류: %s, 결과 : %v", exception.ErrLectureCreditInvalid, err)
		}
	})

	// given
	t.Run("예외 : 요일이 지정되지 않은 경우", func(t *testing.T) {
		// when
		_, err := NewLecture(1006, "알고리즘", 30, 3, Day("토요일"), "11:00", "14:00")
		// then
		if err == nil || err.Error() != exception.ErrLectureDayRequired {
			t.Errorf("기대 : %s, 결과 : %v", exception.ErrLectureDayRequired, err)
		}
	})

	// given
	t.Run("예외 : 강의 종료 시간이 시작 시간보다 빠른 경우", func(t *testing.T) {
		// when
		_, err := NewLecture(1007, "데이터베이스", 30, 3, Tuesday, "16:00", "15:00")
		// then
		if err == nil || err.Error() != exception.ErrLectureTimeOrderInvalid {
			t.Errorf("기대 : %s, 결과 : %v", exception.ErrLectureTimeOrderInvalid, err)
		}
	})
}
