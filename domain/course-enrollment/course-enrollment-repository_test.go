package courseenrollmentdomain

import (
	"context"
	"database/sql"
	"errors"
	"reflect"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestCourseEnrollmentDB_CreateEnrollment(t *testing.T) {
	const (
		studentID = 1
		courseID  = 101
		status    = 1
	)
	constCreateTime := time.Date(2023, 8, 25, 0, 0, 0, 0, time.UTC)
	constUpdateTime := time.Date(2023, 8, 25, 1, 0, 0, 0, time.UTC)

	type fields struct {
		DB *sql.DB
	}
	type args struct {
		ctx              context.Context
		courseEnrollment CourseEnrollment
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    CourseEnrollment
		wantErr bool
	}{
		{
			name: "Success",
			fields: fields{
				DB: func() *sql.DB {
					db, mock, err := sqlmock.New()
					if err != nil {
						t.Fatalf("error creating mock database: %v", err)
					}
					mock.ExpectExec("INSERT INTO course_enrollments").
						WithArgs(studentID, courseID, status, constCreateTime, constUpdateTime).
						WillReturnResult(sqlmock.NewResult(1, 1))
					return db
				}(),
			},
			args: args{
				ctx: context.Background(),
				courseEnrollment: CourseEnrollment{
					StudentID:  studentID,
					CourseID:   courseID,
					Status:     status,
					CreateTime: constCreateTime,
					UpdateTime: constUpdateTime,
				},
			},
			want: CourseEnrollment{
				ID:         1,
				StudentID:  studentID,
				CourseID:   courseID,
				Status:     status,
				CreateTime: constCreateTime,
				UpdateTime: constUpdateTime,
			},
			wantErr: false,
		},
		{
			name: "Error",
			fields: fields{
				DB: func() *sql.DB {
					db, mock, err := sqlmock.New()
					if err != nil {
						t.Fatalf("error creating mock database: %v", err)
					}
					mock.ExpectExec("INSERT INTO course_enrollments").
						WithArgs(studentID, courseID, status, constCreateTime, constUpdateTime).
						WillReturnError(errors.New("insert error"))
					return db
				}(),
			},
			args: args{
				ctx: context.Background(),
				courseEnrollment: CourseEnrollment{
					StudentID:  studentID,
					CourseID:   courseID,
					Status:     status,
					CreateTime: constCreateTime,
					UpdateTime: constUpdateTime,
				},
			},
			want:    CourseEnrollment{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &CourseEnrollmentDB{
				DB: tt.fields.DB,
			}
			got, err := repo.CreateEnrollment(tt.args.ctx, tt.args.courseEnrollment)
			if (err != nil) != tt.wantErr {
				t.Errorf("CourseEnrollmentDB.CreateEnrollment() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CourseEnrollmentDB.CreateEnrollment() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCourseEnrollmentDB_GetEnrollmentByStudentID(t *testing.T) {
	const studentID int64 = 1

	type fields struct {
		DB *sql.DB
	}
	type args struct {
		ctx       context.Context
		studentID int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []CourseEnrollment
		wantErr bool
	}{
		{
			name: "Success",
			fields: fields{
				DB: func() *sql.DB {
					db, mock, err := sqlmock.New()
					if err != nil {
						t.Fatalf("error creating mock database: %v", err)
					}
					rows := sqlmock.NewRows([]string{"id", "student_id", "course_id", "status", "create_time", "update_time"}).
						AddRow(1, studentID, 101, 1, time.Date(2024, 8, 25, 0, 0, 0, 0, time.UTC), time.Date(2024, 8, 25, 0, 0, 0, 0, time.UTC))

					mock.ExpectQuery(`SELECT id, student_id, course_id, status, create_time, update_time FROM course_enrollments WHERE student_id = \? and status = 1`).
						WithArgs(studentID).
						WillReturnRows(rows)
					return db
				}(),
			},
			args: args{
				ctx:       context.Background(),
				studentID: studentID,
			},
			want: []CourseEnrollment{
				{
					ID:         1,
					StudentID:  studentID,
					CourseID:   101,
					Status:     1,
					CreateTime: time.Date(2024, 8, 25, 0, 0, 0, 0, time.UTC),
					UpdateTime: time.Date(2024, 8, 25, 0, 0, 0, 0, time.UTC),
				},
			},
			wantErr: false,
		},
		{
			name: "No Rows",
			fields: fields{
				DB: func() *sql.DB {
					db, mock, err := sqlmock.New()
					if err != nil {
						t.Fatalf("error creating mock database: %v", err)
					}
					mock.ExpectQuery(`SELECT id, student_id, course_id, status, create_time, update_time FROM course_enrollments WHERE student_id = \? and status = 1`).
						WithArgs(studentID).
						WillReturnRows(sqlmock.NewRows([]string{"id", "student_id", "course_id", "status", "create_time", "update_time"}))
					return db
				}(),
			},
			args: args{
				ctx:       context.Background(),
				studentID: studentID,
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "Error",
			fields: fields{
				DB: func() *sql.DB {
					db, mock, err := sqlmock.New()
					if err != nil {
						t.Fatalf("error creating mock database: %v", err)
					}
					mock.ExpectQuery(regexp.QuoteMeta("SELECT id, student_id, course_id, status, create_time, update_time FROM course_enrollments WHERE student_id = ? and status = 1")).
						WithArgs(studentID).
						WillReturnError(errors.New("query error"))
					return db
				}(),
			},
			args: args{
				ctx:       context.Background(),
				studentID: studentID,
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &CourseEnrollmentDB{
				DB: tt.fields.DB,
			}
			got, err := repo.GetEnrollmentByStudentID(tt.args.ctx, tt.args.studentID)
			if (err != nil) != tt.wantErr {
				t.Errorf("CourseEnrollmentDB.GetEnrollmentByStudentID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CourseEnrollmentDB.GetEnrollmentByStudentID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCourseEnrollmentDB_UpdateCourseEnrollmentStatus(t *testing.T) {
	type fields struct {
		DB *sql.DB
	}
	type args struct {
		ctx       context.Context
		studentID int64
		courseID  int64
		newStatus int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Success",
			fields: fields{
				DB: func() *sql.DB {
					db, mock, err := sqlmock.New()
					if err != nil {
						t.Fatalf("error creating mock database: %v", err)
					}
					mock.ExpectExec(`UPDATE course_enrollments SET status = \?, update_time = \? WHERE student_id = \? AND course_id = \?`).
						WithArgs(1, sqlmock.AnyArg(), int64(1), int64(101)).
						WillReturnResult(sqlmock.NewResult(0, 1)) // 1 row affected
					return db
				}(),
			},
			args: args{
				ctx:       context.Background(),
				studentID: 1,
				courseID:  101,
				newStatus: 1,
			},
			wantErr: false,
		},
		{
			name: "No Rows Affected",
			fields: fields{
				DB: func() *sql.DB {
					db, mock, err := sqlmock.New()
					if err != nil {
						t.Fatalf("error creating mock database: %v", err)
					}
					mock.ExpectExec(`UPDATE course_enrollments SET status = \?, update_time = \? WHERE student_id = \? AND course_id = \?`).
						WithArgs(1, sqlmock.AnyArg(), int64(1), int64(101)).
						WillReturnResult(sqlmock.NewResult(0, 0)) // No rows affected
					return db
				}(),
			},
			args: args{
				ctx:       context.Background(),
				studentID: 1,
				courseID:  101,
				newStatus: 1,
			},
			wantErr: true,
		},
		{
			name: "Query Error",
			fields: fields{
				DB: func() *sql.DB {
					db, mock, err := sqlmock.New()
					if err != nil {
						t.Fatalf("error creating mock database: %v", err)
					}
					mock.ExpectExec(`UPDATE course_enrollments SET status = \?, update_time = \? WHERE student_id = \? AND course_id = \?`).
						WithArgs(1, sqlmock.AnyArg(), int64(1), int64(101)).
						WillReturnError(errors.New("update failed"))
					return db
				}(),
			},
			args: args{
				ctx:       context.Background(),
				studentID: 1,
				courseID:  101,
				newStatus: 1,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &CourseEnrollmentDB{
				DB: tt.fields.DB,
			}
			err := repo.UpdateCourseEnrollmentStatus(tt.args.ctx, tt.args.studentID, tt.args.courseID, tt.args.newStatus)
			if (err != nil) != tt.wantErr {
				t.Errorf("CourseEnrollmentDB.UpdateCourseEnrollmentStatus() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCourseEnrollmentDB_GetListClassmates(t *testing.T) {

	timestamp := time.Date(2024, time.August, 1, 0, 0, 0, 0, time.UTC)

	type fields struct {
		DB *sql.DB
	}
	type args struct {
		ctx       context.Context
		studentID int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []CourseEnrollment
		wantErr bool
	}{
		{
			name: "Success",
			fields: fields{
				DB: func() *sql.DB {
					db, mock, err := sqlmock.New()
					if err != nil {
						t.Fatalf("error creating mock database: %v", err)
					}
					rows := sqlmock.NewRows([]string{"id", "student_id", "course_id", "status", "create_time", "update_time"}).
						AddRow(1, 101, 1001, 1, timestamp, timestamp).
						AddRow(2, 102, 1001, 1, timestamp, timestamp)
					mock.ExpectQuery(`SELECT ce.id, ce.student_id, ce.course_id, ce.status, ce.create_time, ce.update_time FROM course_enrollments ce JOIN course_enrollments ce2 ON ce.course_id = ce2.course_id WHERE ce2.student_id = \? AND ce.student_id != \? and ce2.status = 1 and ce.status = 1`).
						WithArgs(int64(1), int64(1)).
						WillReturnRows(rows)
					return db
				}(),
			},
			args: args{
				ctx:       context.Background(),
				studentID: 1,
			},
			want: []CourseEnrollment{
				{ID: 1, StudentID: 101, CourseID: 1001, Status: 1, CreateTime: timestamp, UpdateTime: timestamp},
				{ID: 2, StudentID: 102, CourseID: 1001, Status: 1, CreateTime: timestamp, UpdateTime: timestamp},
			},
			wantErr: false,
		},
		{
			name: "Query Error",
			fields: fields{
				DB: func() *sql.DB {
					db, mock, err := sqlmock.New()
					if err != nil {
						t.Fatalf("error creating mock database: %v", err)
					}
					mock.ExpectQuery(`SELECT ce.id, ce.student_id, ce.course_id, ce.status, ce.create_time, ce.update_time FROM course_enrollments ce JOIN course_enrollments ce2 ON ce.course_id = ce2.course_id WHERE ce2.student_id = \? AND ce.student_id != \? and ce2.status = 1 and ce.status = 1`).
						WithArgs(int64(1), int64(1)).
						WillReturnError(sql.ErrConnDone)
					return db
				}(),
			},
			args: args{
				ctx:       context.Background(),
				studentID: 1,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Scan Error",
			fields: fields{
				DB: func() *sql.DB {
					db, mock, err := sqlmock.New()
					if err != nil {
						t.Fatalf("error creating mock database: %v", err)
					}
					rows := sqlmock.NewRows([]string{"id", "student_id", "course_id", "status", "create_time", "update_time"}).
						AddRow("invalid", 101, 1001, 1, timestamp, timestamp)
					mock.ExpectQuery(`SELECT ce.id, ce.student_id, ce.course_id, ce.status, ce.create_time, ce.update_time FROM course_enrollments ce JOIN course_enrollments ce2 ON ce.course_id = ce2.course_id WHERE ce2.student_id = \? AND ce.student_id != \? and ce2.status = 1 and ce.status = 1`).
						WithArgs(int64(1), int64(1)).
						WillReturnRows(rows)
					return db
				}(),
			},
			args: args{
				ctx:       context.Background(),
				studentID: 1,
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &CourseEnrollmentDB{
				DB: tt.fields.DB,
			}
			got, err := repo.GetListClassmates(tt.args.ctx, tt.args.studentID)
			if (err != nil) != tt.wantErr {
				t.Errorf("CourseEnrollmentDB.GetListClassmates() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CourseEnrollmentDB.GetListClassmates() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCourseEnrollmentDB_GetEnrollmentByStudentIDAndCourseID(t *testing.T) {
	const (
		studentID = 1
		courseID  = 101
		status    = 1
	)
	timestamp := time.Date(2023, 8, 25, 0, 0, 0, 0, time.UTC)

	type fields struct {
		DB *sql.DB
	}
	type args struct {
		ctx       context.Context
		studentID int64
		courseID  int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []CourseEnrollment
		wantErr bool
	}{
		{
			name: "Success",
			fields: fields{
				DB: func() *sql.DB {
					db, mock, err := sqlmock.New()
					if err != nil {
						t.Fatalf("error creating mock database: %v", err)
					}
					rows := sqlmock.NewRows([]string{"id", "student_id", "course_id", "status", "create_time", "update_time"}).
						AddRow(1, studentID, courseID, status, timestamp, timestamp)
					mock.ExpectQuery(`SELECT id, student_id, course_id, status, create_time, update_time FROM course_enrollments WHERE student_id = \? AND course_id = \?`).
						WithArgs(studentID, courseID).
						WillReturnRows(rows)
					return db
				}(),
			},
			args: args{
				ctx:       context.Background(),
				studentID: studentID,
				courseID:  courseID,
			},
			want: []CourseEnrollment{
				{
					ID:         1,
					StudentID:  studentID,
					CourseID:   courseID,
					Status:     status,
					CreateTime: timestamp,
					UpdateTime: timestamp,
				},
			},
			wantErr: false,
		},
		{
			name: "Error",
			fields: fields{
				DB: func() *sql.DB {
					db, mock, err := sqlmock.New()
					if err != nil {
						t.Fatalf("error creating mock database: %v", err)
					}
					mock.ExpectQuery("SELECT id, student_id, course_id, status, create_time, update_time FROM course_enrollments WHERE student_id = ? AND course_id = ?").
						WithArgs(studentID, courseID).
						WillReturnError(errors.New("query error"))
					return db
				}(),
			},
			args: args{
				ctx:       context.Background(),
				studentID: studentID,
				courseID:  courseID,
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &CourseEnrollmentDB{
				DB: tt.fields.DB,
			}
			got, err := repo.GetEnrollmentByStudentIDAndCourseID(tt.args.ctx, tt.args.studentID, tt.args.courseID)
			if (err != nil) != tt.wantErr {
				t.Errorf("CourseEnrollmentDB.GetEnrollmentByStudentIDAndCourseID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CourseEnrollmentDB.GetEnrollmentByStudentIDAndCourseID() = %v, want %v", got, tt.want)
			}
		})
	}
}
