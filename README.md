# Golang : 수강신청 시스템 

## 1. 소개

Go 언어와 Echo 프레임워크를 사용하여 구현한 강좌 수강신청 웹 애플리케이션입니다. 관리자는 강좌를 등록하고 관리할 수 있으며, 학생은 강좌를 조회하고 수강신청할 수 있습니다.

빠른 수강신청 프로세스를 위해 인증/인가 프로세스는 간단하게 학번으로만 검증하도록 하였으며, 수강 신청 및 수강 신청 취소 시 동시성 제어를를 위한 메모리 기반 락을 구현하여 단일 서버 환경에서의 경쟁 조건을 방지하고자 하였습니다.
## 2. 기능 요구사항

### -1. 관리자 기능

#### 강좌 등록
- **강좌번호**: 1000~9999 사이의 숫자
- **강좌명**: 2~20자 사이의 문자열 (중복 불가)
- **정원**: 1명 이상 30명 이하
- **학점**: 1학점 이상 6학점 이하
- **요일**: 월요일~금요일 중 선택
- **시간**: 시작시간과 종료시간 입력 (HH:MM 형식)
- **검증**: 강좌명 및 강좌번호 중복 체크, 시간 형식 및 유효성 검증

#### 강좌 조회
- 등록된 모든 강좌 목록 조회
- 각 강좌의 현재 수강 인원 및 정원 표시

#### 강좌 삭제
- 등록된 강좌 삭제
- 외래키 제약조건(`ON DELETE CASCADE`)으로 관련 수강신청 삭제

### -2. 학생 기능

#### 학생 등록
- **학번**: 1000~9999 사이의 4자리 숫자

#### 강좌 목록 조회
- 등록된 모든 강좌 목록 조회
- 각 강좌의 학점, 현재 수강 인원, 정원, 요일, 시간 정보 표시

#### 수강신청
- 강좌별 수강신청 버튼을 통한 신청
- **검증 항목**:
  - 학생 존재 여부 확인
  - 강좌 존재 여부 확인
  - 정원 초과 여부 확인
  - 시간 충돌 검사 (같은 요일에 겹치는 시간 방지)
  - 총 학점 제한 (18학점 초과 불가)

#### 수강신청 내역 조회
- 본인이 신청한 강좌 목록 조회
- 각 강좌의 상세 정보 표시

#### 수강신청 취소
- 동시성 제어 락 획득 후, 수강신청 내역 삭제

### -3. 웹 페이지

- **메인 페이지** (`/`): 시스템 소개 및 학번 입력을 통한 대시보드 접속
- **관리자 대시보드** (`/admin/dashboard`): 강좌 등록, 조회, 삭제 기능
- **학생 대시보드** (`/client/dashboard?studentId`): 강좌 조회 및 수강신청 기능

## 3. 프로젝트 구조

```
golang-course-registration/
├── common/
│   └── exception/           # 예외 메시지 정의
├── config/                  # 설정 관리 (환경 변수)
├── controller/
│   ├── api/                 # REST API 컨트롤러
│   │   ├── admin_controller.go
│   │   └── client_controller.go
│   ├── dto/                 # 데이터 전송 객체
│   │   ├── lecture_dto.go
│   │   └── enrollment_dto.go
│   └── web/                 # 웹 페이지 컨트롤러
│       └── page_controller.go
├── infrastructure/
│   ├── database/            # 데이터베이스 연결 (Supabase)
│   └── server/              # 서버 설정 및 라우팅
├── model/                   # 도메인 모델
│   ├── student.go
│   ├── student_test.go
│   ├── lecture.go
│   ├── lecture_test.go
│   ├── enrollment.go
│   ├── enrollment_test.go
│   └── day.go
├── repository/              # 데이터 접근 계층
│   ├── student_repository.go
│   ├── lecture_repository.go
│   └── enrollment_repository.go
├── service/                 # 비즈니스 로직 계층
│   ├── student_service.go
│   ├── student_service_test.go
│   ├── lecture_service.go
│   ├── lecture_service_test.go
│   ├── enrollment_service.go
│   └── enrollment_service_test.go
│
├── view/                    # HTML 템플릿 및 정적 파일
│   ├── templates/           # HTML 템플릿
│   │   ├── base.html
│   │   ├── index.html
│   │   ├── admin.html
│   │   └── client.html
│   ├── style/               # CSS 스타일 파일
│   │   ├── admin_styles.html
│   │   └── client_styles.html
│   └── script/              # JavaScript 파일
│       ├── admin_scripts.html
│       ├── client_scripts.html
│       └── index_scripts.html
├── main.go                  # 애플리케이션 진입점
├── go.mod                   # Go 모듈 정의
└── go.sum                   # 의존성 체크섬
```

## 4. 기술 스택

- **언어**: Go 1.24
- **웹 프레임워크**: Echo v4
- **데이터베이스**: Supabase (PostgreSQL)
- **템플릿 엔진**: Go html/template
- **환경 변수 관리**: godotenv
- **아키텍처**: 계층형 아키텍처 (Controller → Service → Repository → Database)

## 5. 주요 구현 기능

### - 5.1 동시성 제어

#### 메모리 기반 락 (단일 서버 환경용)
- 강좌별 개별 락(`sync.Mutex`) 관리
- 수강신청 시 해당 강좌의 락을 획득하여 동시성 제어

### - 5.2 학점 관리

#### 총 학점 제한 (18학점)
- 수강신청 시 기존 수강신청 강좌들의 학점 합산
- 새 강좌 학점 추가 시 18학점 초과 여부 검증
- 강좌별 학점은 1~6학점 범위

### - 5.3 시간 중복 검사

#### 같은 요일 시간 중복 방지
- 수강신청 시 기존 수강신청 강좌들과 시간 비교
- 같은 요일에서 시간이 겹치는 경우 수강신청 불가

### - 5.4 강좌 삭제 시, 데이터 일관성 보장

#### - CASCADE 삭제
- 강좌 삭제 시 관련 수강신청 삭제
- 학생 삭제 시 관련 수강신청 삭제

### - 5.5 입력 검증

#### 강좌 등록 검증
- 강좌번호: 1000~9999
- 강좌명: 2~20자
- 정원: 1~30명
- 학점: 1~6학점
- 시간 형식: HH:MM
- 종료 시간 > 시작 시간

#### 학생 등록 검증
- 학번: 1000~9999

## 6. 예외 처리

### 6.1 예외 메시지 정의

모든 예외 메시지는 `common/exception/messages.go`에 관리

- Student 관련 예외
- Lecture 관련 예외
- Enrollment 관련 예외
- Controller 관련 예외

## 7. 실행 및 배포

### 7.1 로컬 환경에서 실행

#### 1. 환경 변수 설정
`.env` 파일을 생성하고 다음 변수를 설정.

```
SUPABASE_URL=your_supabase_url
SUPABASE_KEY=your_supabase_key
PORT=8080
```

#### 2. 의존성 설치 및 실행
```bash
go mod download
go run main.go
```

### 7.2 Docker를 이용한 배포

#### 1. Docker 이미지 빌드
```bash
docker build -t golang-course-system .
```

#### 2. Docker 컨테이너 실행
환경 변수를 `-e` 옵션으로 주입하여 컨테이너를 실행.

```bash
docker run -p 8080:8080 \
  -e SUPABASE_URL="your_supabase_url" \
  -e SUPABASE_ANON_KEY="your_supabase_anon_key" \
  -e PORT="8080" \
  golang-course-system
```

## 8. 테스트

### 테스트 실행 방법

```bash
go test -v ./...
```

### 주요 테스트 케이스
- **도메인 모델 (`/model`)**: 각 도메인 객체(`Lecture`, `Student`, `Enrollment`) 유효성 검증 및 비즈니스 로직을 테스트합니다.
- **서비스 계층 (`/service`)**:
  - **LectureService**: 강의 생성, 조회, 삭제 기능 및 중복 처리와 같은 예외 상황을 검증합니다.
  - **StudentService**: 학생 등록 및 유효성 검증을 테스트합니다.
  - **EnrollmentService**: 수강 신청 및 취소 로직을 검증하며, 정원 초과, 시간 충돌, 학점 제한 등 다양한 예외 케이스를 포함합니다.

## 9. API 엔드포인트

### 관리자 API
- `POST /api/v1/admin/lectures`: 강좌 등록
- `GET /api/v1/admin/lectures`: 강좌 목록 조회
- `DELETE /api/v1/admin/lectures/:id`: 강좌 삭제

### 학생 API

- `POST /api/v1/client/students`: 학생 등록
- `GET /api/v1/client/lectures`: 강좌 목록 조회
- `POST /api/v1/client/enrollments`: 수강신청
- `GET /api/v1/client/enrollments/:studentId`: 수강신청 내역 조회
- `DELETE /api/v1/client/enrollments/:studentId/:lectureId`: 수강신청 취소

## 9. DB 스키마 

```postgresql
CREATE TABLE enrollments (
  id bigint GENERATED ALWAYS AS IDENTITY NOT NULL,
  student_id bigint NOT NULL,
  lecture_id bigint NOT NULL,
  CONSTRAINT enrollments_pkey PRIMARY KEY (id),
  CONSTRAINT enrollments_lecture_id_fkey FOREIGN KEY (lecture_id) REFERENCES lectures(id) ON DELETE CASCADE,
  CONSTRAINT enrollments_student_id_fkey FOREIGN KEY (student_id) REFERENCES students(id) ON DELETE CASCADE
);

CREATE TABLE lectures (
  id bigint GENERATED ALWAYS AS IDENTITY NOT NULL UNIQUE,
  name character varying NOT NULL,
  capacity bigint NOT NULL,
  day character varying NOT NULL,
  start_time character varying NOT NULL,
  end_time character varying NOT NULL,
  current_enrollment bigint NOT NULL DEFAULT 0,
  credit bigint NOT NULL,
  CONSTRAINT lectures_pkey PRIMARY KEY (id)
);

CREATE TABLE students (
  id bigint GENERATED ALWAYS AS IDENTITY NOT NULL,
  CONSTRAINT students_pkey PRIMARY KEY (id)
);
```