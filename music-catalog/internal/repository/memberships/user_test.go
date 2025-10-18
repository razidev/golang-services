package memberships

import (
	"music-catalog/internal/models/memberships"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Test_repository_CreateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDb, err := gorm.Open(postgres.New(postgres.Config{Conn: db}))
	assert.NoError(t, err)

	type args struct {
		model memberships.User
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		mockFn  func(args args)
	}{
		{
			name: "Succes",
			args: args{
				model: memberships.User{
					Email:     "test@gmail.com",
					Username:  "username1",
					Password:  "password",
					CreatedBy: "test@gmail.com",
					UpdateBy:  "test@gmail.com",
				},
			},
			wantErr: false,
			mockFn: func(args args) {
				mock.ExpectBegin()

				mock.ExpectQuery(`INSERT INTO "users" (.+) VALUES (.+)`).
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						args.model.Email,
						args.model.Username,
						args.model.Password,
						args.model.CreatedBy,
						args.model.UpdateBy,
					).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

				mock.ExpectCommit()
			},
		},
		{
			name: "error",
			args: args{
				model: memberships.User{
					Email:     "test@gmail.com",
					Username:  "username1",
					Password:  "password",
					CreatedBy: "test@gmail.com",
					UpdateBy:  "test@gmail.com",
				},
			},
			wantErr: true,
			mockFn: func(args args) {
				mock.ExpectBegin()

				mock.ExpectQuery(`INSERT INTO "users" (.+) VALUES (.+)`).
					WithArgs(
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						args.model.Email,
						args.model.Username,
						args.model.Password,
						args.model.CreatedBy,
						args.model.UpdateBy,
					).
					WillReturnError(assert.AnError)

				mock.ExpectRollback()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn(tt.args)
			r := NewRepository(gormDb)
			gotErr := r.CreateUser(tt.args.model)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("CreateUser() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("CreateUser() succeeded unexpectedly")
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func Test_repository_GetUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	gormDb, err := gorm.Open(postgres.New(postgres.Config{Conn: db}))
	assert.NoError(t, err)

	now := time.Now()

	type args struct {
		email    string
		username string
		id       uint
	}

	tests := []struct {
		name    string
		args    args
		want    *memberships.User
		wantErr bool
		mockFn  func(args args)
	}{
		{
			name: "success",
			args: args{
				email:    "test@gmail.com",
				username: "username1",
			},
			want: &memberships.User{
				Model: gorm.Model{
					ID:        1,
					CreatedAt: now,
					UpdatedAt: now,
				},
				Email:     "test@gmail.com",
				Username:  "username1",
				Password:  "password1",
				CreatedBy: "test@gmail.com",
				UpdateBy:  "test@gmail.com",
			},
			wantErr: false,
			mockFn: func(args args) {
				mock.ExpectQuery(`SELECT \* FROM "users" .+`).WithArgs(args.email, args.username, args.id, sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "email", "username",
						"password", "created_by", "update_by"}).
						AddRow(1, now, now, "test@gmail.com", "username1", "password1", "test@gmail.com", "test@gmail.com"))
			},
		},
		{
			name: "failed",
			args: args{
				email:    "test@gmail.com",
				username: "username1",
			},
			want:    nil,
			wantErr: true,
			mockFn: func(args args) {
				mock.ExpectQuery(`SELECT \* FROM "users" .+`).WithArgs(args.email, args.username, args.id, sqlmock.AnyArg()).
					WillReturnError(assert.AnError)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn(tt.args)

			r := NewRepository(gormDb)
			got, gotErr := r.GetUser(tt.args.email, tt.args.username, tt.args.id)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("GetUser() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("GetUser() succeeded unexpectedly")
			}
			// TODO: update the condition below to compare got with tt.want.
			// Check the UpdateBy field specifically (was missing in previous mock results)
			if got.UpdateBy != tt.want.UpdateBy {
				t.Errorf("GetUser().UpdateBy = %v, want %v", got.UpdateBy, tt.want.UpdateBy)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
