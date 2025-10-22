package memberships

import (
	"music-catalog/internal/configs"
	"music-catalog/internal/models/memberships"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

func Test_service_SignUp(t *testing.T) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()

	mockRepo := NewMockrepository(ctrlMock)

	type args struct {
		request memberships.SignUpRequest
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
				request: memberships.SignUpRequest{
					Email:    "test@gmail.com",
					Username: "testUsername",
					Password: "password123",
				},
			},
			wantErr: false,
			mockFn: func(args args) {
				mockRepo.EXPECT().GetUser(args.request.Email, args.request.Username, uint(0)).Return(nil, gorm.ErrRecordNotFound)
				mockRepo.EXPECT().CreateUser(gomock.Any()).Return(nil)
			},
		},
		{
			name: "Failed to get user",
			args: args{
				request: memberships.SignUpRequest{
					Email:    "test@gmail.com",
					Username: "testUsername",
					Password: "password123",
				},
			},
			wantErr: true,
			mockFn: func(args args) {
				mockRepo.EXPECT().GetUser(args.request.Email, args.request.Username, uint(0)).Return(nil, assert.AnError)
			},
		},
		{
			name: "Failed when create user",
			args: args{
				request: memberships.SignUpRequest{
					Email:    "test@gmail.com",
					Username: "testUsername",
					Password: "password123",
				},
			},
			wantErr: true,
			mockFn: func(args args) {
				mockRepo.EXPECT().GetUser(args.request.Email, args.request.Username, uint(0)).Return(nil, gorm.ErrRecordNotFound)
				mockRepo.EXPECT().CreateUser(gomock.Any()).Return(assert.AnError)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn(tt.args)
			s := NewService(&configs.Config{}, mockRepo)
			gotErr := s.SignUp(tt.args.request)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("SignUp() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("SignUp() succeeded unexpectedly")
			}
		})
	}
}
