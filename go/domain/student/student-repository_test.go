package studentdomain

import (
	"context"
	"database/sql"
	"reflect"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestStudentDB_GetStudentByID(t *testing.T) {
	const (
		studentID    = 1
		studentEmail = "test@example.com"
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
		want    *Student
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
					rows := sqlmock.NewRows([]string{"id", "email", "create_time", "update_time"}).
						AddRow(studentID, studentEmail, constCreateTime, constUpdateTime)
					mock.ExpectQuery("SELECT id, email, create_time, update_time FROM students WHERE id = ?").
						WithArgs(studentID).
						WillReturnRows(rows)
					return db
				}(),
			},
			args: args{
				ctx: context.Background(),
				id:  studentID,
			},
			want: &Student{
				ID:         studentID,
				Email:      studentEmail,
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
					mock.ExpectQuery("SELECT id, email, create_time, update_time FROM students WHERE id = ?").
						WithArgs(studentID).
						WillReturnRows(sqlmock.NewRows([]string{"id", "email", "create_time", "update_time"}))
					return db
				}(),
			},
			args: args{
				ctx: context.Background(),
				id:  studentID,
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
					mock.ExpectQuery("SELECT id, email, create_time, update_time FROM students WHERE id = ?").
						WithArgs(studentID).
						WillReturnError(sql.ErrConnDone)
					return db
				}(),
			},
			args: args{
				ctx: context.Background(),
				id:  studentID,
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &StudentDB{
				DB: tt.fields.DB,
			}
			got, err := repo.GetStudentByID(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("StudentDB.GetStudentByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StudentDB.GetStudentByID() = %v, want %v", got, tt.want)
			}
		})
	}
}
