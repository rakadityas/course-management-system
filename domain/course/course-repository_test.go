package coursedomain

import (
	"context"
	"database/sql"
	"reflect"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestCourseDB_GetCourseByID(t *testing.T) {
	const (
		courseID   = 1
		courseName = "Introduction to Go"
	)

	constCreateTime := time.Date(2023, 8, 25, 0, 0, 0, 0, time.UTC)
	constUpdateTime := time.Date(2023, 8, 25, 1, 0, 0, 0, time.UTC)

	type fields struct {
		DB *sql.DB
	}
	type args struct {
		ctx context.Context
		id  int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Course
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
					rows := sqlmock.NewRows([]string{"id", "name", "create_time", "update_time"}).
						AddRow(courseID, courseName, constCreateTime, constUpdateTime)
					mock.ExpectQuery("SELECT id, name, create_time, update_time FROM courses WHERE id = ?").
						WithArgs(courseID).
						WillReturnRows(rows)
					return db
				}(),
			},
			args: args{
				ctx: context.Background(),
				id:  courseID,
			},
			want: &Course{
				ID:         courseID,
				Name:       courseName,
				CreateTime: constCreateTime,
				UpdateTime: constUpdateTime,
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
					mock.ExpectQuery("SELECT id, name, create_time, update_time FROM courses WHERE id = ?").
						WithArgs(courseID).
						WillReturnRows(sqlmock.NewRows([]string{"id", "name", "create_time", "update_time"}))
					return db
				}(),
			},
			args: args{
				ctx: context.Background(),
				id:  courseID,
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
					mock.ExpectQuery("SELECT id, name, create_time, update_time FROM courses WHERE id = ?").
						WithArgs(courseID).
						WillReturnError(sql.ErrConnDone)
					return db
				}(),
			},
			args: args{
				ctx: context.Background(),
				id:  courseID,
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &CourseDB{
				DB: tt.fields.DB,
			}
			got, err := repo.GetCourseByID(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("CourseDB.GetCourseByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CourseDB.GetCourseByID() = %v, want %v", got, tt.want)
			}
		})
	}
}
