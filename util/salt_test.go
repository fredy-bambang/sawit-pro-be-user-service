package util

import (
	"testing"
)

func TestGenerateSalt(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "Valid Salt",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			saltResult := GenerateSalt()
			if saltResult == "" {
				t.Errorf("Generated salt should not return empty string")
				return
			}
		})
	}
}

func TestHashPassword(t *testing.T) {
	type args struct {
		password string
		salt     string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Valid Hash Password",
			args: args{
				password: "password",
				salt:     "salt",
			},
			want:    "S8D9UH6TpgB2gCE0HscmxXwAy1WkcCoWUBMTZVAM9HE=",
			wantErr: false,
		},
		{
			name: "Not Valid Password",
			args: args{
				password: "password",
				salt:     "salt",
			},
			want:    "invalid password",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := HashPassword(tt.args.password, tt.args.salt)
			if tt.wantErr && got != tt.want {
				return
			}
			if got != tt.want {
				t.Errorf("HashPassword() = %v, want %v", got, tt.want)
			}
		})
	}
}
