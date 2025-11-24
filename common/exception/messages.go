package exception

// Student 관련 예외 메시지
const (
	ErrStudentIDInvalid = "학번(ID)은 1000 ~ 9999 사이의 숫자여야 합니다"
)

// Lecture 관련 예외 메시지
const (
	ErrLectureNameRequired     = "강좌명은 2~20자 사이여야 합니다"
	ErrLectureIDInvalid        = "강좌번호는 1000 ~ 9999 사이의 숫자여야 합니다"
	ErrLectureCapacityInvalid  = "정원은 1명 이상, 30명 이하여야 합니다"
	ErrLectureDayRequired      = "강좌 요일은 필수입니다"
	ErrLectureTimeRequired     = "시작/종료 시간은 필수입니다"
	ErrLectureTimeOrderInvalid = "종료 시간은 시작 시간 이후여야 합니다"
	ErrLectureNameDuplicate    = "이미 존재하는 강좌명입니다"
	ErrLectureIDDuplicate      = "이미 존재하는 강좌번호입니다"
	ErrLectureCreditInvalid    = "학점은 1학점 이상, 6학점 이하여야 합니다"
)

// Enrollment 관련 예외 메시지
const (
	ErrEnrollmentLectureIDRequired = "강좌번호는 필수입니다"
	ErrStudentNotFound             = "존재하지 않는 학생입니다"
	ErrLectureNotFound             = "존재하지 않는 강좌입니다"
	ErrTimeConflict                = "강좌와 시간이 중복됩니다"
	ErrLectureCapacityExceeded     = "강좌 정원이 초과되었습니다"
	ErrCreditLimitExceeded         = "총 학점이 18학점을 초과할 수 없습니다"
)

// Controller 관련 예외 메시지
const (
	ErrInvalidRequestBody = "요청 본문이 올바르지 않습니다"
	ErrLectureListFailed  = "강좌 목록 조회 실패"
	ErrStudentIDNotNumber = "학번은 숫자여야 합니다"
)

// TimeConflictMessage 시간 충돌 메시지 생성
func TimeConflictMessage(lectureName string) string {
	return lectureName + " " + ErrTimeConflict
}
