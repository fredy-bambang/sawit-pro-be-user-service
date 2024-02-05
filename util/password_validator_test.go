package util

import "testing"

func TestValidatePassword(t *testing.T) {
	type args struct {
		password string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Valid Password",
			args: args{
				password: "Password1!",
			},
			want: true,
		},
		{
			name: "Not Valid Password, need at least 1 uppercase, 1 number, and 1 special character",
			args: args{
				password: "password",
			},
			want: false,
		},
		{
			name: "Not Valid Password, need at least 1 number and 1 special character",
			args: args{
				password: "Password",
			},
			want: false,
		},
		{
			name: "Not Valid Password, need at least 1 special character",
			args: args{
				password: "Password1",
			},
			want: false,
		},
		{
			name: "Not Valid Password, need at least 6 char",
			args: args{
				password: "Passw",
			},
			want: false,
		},
		{
			name: "Not Valid Password, max password length is 64 char",
			args: args{
				password: "1234567890123456789012345678901234567890123456789012345678901234567890",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidatePassword(tt.args.password); got != tt.want {
				t.Errorf("ValidatePassword() = %v, want %v", got, tt.want)
			}
		})
	}
}
