package token

import "testing"

func TestGetToken(t *testing.T) {
	tests := []struct {
		name           string
		authentication string
		want           string
	}{
		{
			name:           "valid token",
			authentication: "Bearer token",
			want:           "token",
		},
		{
			name:           "invalid token",
			authentication: "Bearer",
			want:           "",
		},
		{
			name:           "invalid token",
			authentication: "Bearer token token",
			want:           "",
		},
		{
			name:           "invalid token",
			authentication: "",
			want:           "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetToken(tt.authentication); got != tt.want {
				t.Errorf("GetToken() = %v, want %v", got, tt.want)
			}
		})
	}
}
