package memberships

import (
	"music-catalog/internal/configs"
	"music-catalog/internal/models/memberships"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

func Test_service_Login(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	mockRepo := NewMockrepository(ctrlMock)
	type args struct {
		request memberships.LoginRequest
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
		mockFn  func(args args)
	}{
		{
			name: "Success",
			args: args{
				request: memberships.LoginRequest{
					Email:    "razisyahputro@gmail.com",
					Password: "test123",
				},
			},
			wantErr: false,
			mockFn: func(args args) {
				mockRepo.EXPECT().GetUser(args.request.Email, "", uint(0)).Return(&memberships.User{
					Model: gorm.Model{
						ID: 1,
					},
					Email:    "razisyahputro@gmail.com",
					Username: "Razi",
					Password: "$2a$10$rhI2cbzqyKjACLl1OtqBT.Og7N1suW5RdzsQeE7Fxosgqy0X6H00K",
				}, nil)
			},
		},
		{
			name: "Failde when get user",
			args: args{
				request: memberships.LoginRequest{
					Email:    "razisyahputro@gmail.com",
					Password: "test123",
				},
			},
			wantErr: true,
			mockFn: func(args args) {
				mockRepo.EXPECT().GetUser(args.request.Email, "", uint(0)).Return(nil, assert.AnError)
			},
		},
		{
			name: "Failed - wrong password",
			args: args{
				request: memberships.LoginRequest{
					Email:    "razisyahputro@gmail.com",
					Password: "test123",
				},
			},
			wantErr: true,
			mockFn: func(args args) {
				mockRepo.EXPECT().GetUser(args.request.Email, "", uint(0)).Return(&memberships.User{
					Model: gorm.Model{
						ID: 1,
					},
					Email:    "razisyahputro@gmail.com",
					Username: "Razi",
					Password: "Wrong password",
				}, nil)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn(tt.args)
			s := NewService(&configs.Config{
				Service: configs.Service{
					SecretJWT: "tokenjwt",
				},
			}, mockRepo)
			got, gotErr := s.Login(tt.args.request)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("Login() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("Login() succeeded unexpectedly")
			}

			assert.NotEmpty(t, got)
		})
	}
}
