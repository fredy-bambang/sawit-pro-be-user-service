package repository

import (
	"errors"
	"testing"

	"github.com/SawitProRecruitment/UserService/models"
	"github.com/SawitProRecruitment/UserService/repository/mocks"
	"github.com/stretchr/testify/assert"
)

type TestingT interface {
	Errorf(format string, args ...interface{})
}

func TestPgUserRepository_Create(t *testing.T) {
	type args struct {
		user *models.User
	}
	tests := []struct {
		name    string
		args    args
		want    *models.User
		wantErr bool
	}{
		{
			name: "Success Create",
			args: args{
				&models.User{
					PhoneNumber: "+6281234567890",
					Fullname:    "John Doe",
					Password:    "A1234*",
				},
			},
			wantErr: false,
		},
		{
			name: "Fail Create",
			args: args{
				&models.User{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mocks.UserRepository{}
			if !tt.wantErr {
				repo.On("Create", tt.args.user).Return(nil)
			} else {
				repo.On("Create", tt.args.user).Return(errors.New("Fail create user"))
			}

			err := repo.Create(tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("PgUserRepository.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestPgUserRepository_FindByID(t *testing.T) {
	type args struct {
		id int
	}
	tests := []struct {
		name    string
		args    args
		want    *models.User
		wantErr bool
	}{
		{
			name: "Success FindByID",
			args: args{
				1,
			},
			want: &models.User{
				ID:          1,
				PhoneNumber: "+6281234567890",
				Fullname:    "John Doe",
				Password:    "A1234*",
				SaltToken:   "1234",
			},
			wantErr: false,
		},
		{
			name: "Success FindByID",
			args: args{
				100,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mocks.UserRepository{}
			if !tt.wantErr {
				repo.On("FindByID", tt.args.id).Return(&models.User{
					ID:          1,
					PhoneNumber: "+6281234567890",
					Fullname:    "John Doe",
					Password:    "A1234*",
					SaltToken:   "1234",
				}, nil)
			} else {
				repo.On("FindByID", tt.args.id).Return(nil, errors.New("Record Not Found"))
			}

			user, err := repo.FindByID(tt.args.id)
			if !assert.Equal(t, tt.want, user) {
				t.Errorf("PgUserRepository.FindByID() = %v, want %v", user, tt.want)
				return
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("PgUserRepository.FindByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestPgUserRepository_Update(t *testing.T) {
	type args struct {
		user *models.User
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Success Update",
			args: args{
				&models.User{
					ID:          1,
					PhoneNumber: "+6281234567890",
					Fullname:    "John Doe Update",
					Password:    "A1234*",
					SaltToken:   "1234",
				},
			},
			wantErr: false,
		},
		{
			name: "Fail Update",
			args: args{
				&models.User{
					ID:          1,
					PhoneNumber: "",
					Fullname:    "John Doe Update",
					Password:    "A1234*",
					SaltToken:   "1234",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mocks.UserRepository{}
			if !tt.wantErr {
				repo.On("Update", tt.args.user).Return(nil)
			} else {
				repo.On("Update", tt.args.user).Return(errors.New("Fail To Update"))
			}

			err := repo.Update(tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("PgUserRepository.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestPgUserRepository_FindByPhone(t *testing.T) {
	type args struct {
		phone string
	}
	tests := []struct {
		name    string
		args    args
		want    *models.User
		wantErr bool
	}{
		{
			name: "Success FindByPhone",
			args: args{
				"+6281234567890",
			},
			want: &models.User{
				ID:          1,
				PhoneNumber: "+6281234567890",
				Fullname:    "John Doe",
				Password:    "A1234*",
				SaltToken:   "1234",
			},
			wantErr: false,
		},
		{
			name: "Success FindByPhone",
			args: args{
				"+62812345678111",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mocks.UserRepository{}
			if !tt.wantErr {
				repo.On("FindByPhone", tt.args.phone).Return(&models.User{
					ID:          1,
					PhoneNumber: "+6281234567890",
					Fullname:    "John Doe",
					Password:    "A1234*",
					SaltToken:   "1234",
				}, nil)
			} else {
				repo.On("FindByPhone", tt.args.phone).Return(nil, errors.New("Fail Find By Phone"))
			}

			user, err := repo.FindByPhone(tt.args.phone)
			if !assert.Equal(t, tt.want, user) {
				t.Errorf("PgUserRepository.FindByPhone() = %v, want %v", user, tt.want)
				return
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("PgUserRepository.FindByPhone() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
