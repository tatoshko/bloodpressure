package core

import (
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func TestGetUserName(t *testing.T) {
	tests := []struct {
		name    string
		user    *tgbotapi.User
		mention bool
		want    string
	}{
		{
			name: "Username only - no mention",
			user: &tgbotapi.User{
				UserName:  "johndoe",
				FirstName: "John",
				LastName:  "Doe",
			},
			mention: false,
			want:    "johndoe",
		},
		{
			name: "Username only - with mention",
			user: &tgbotapi.User{
				UserName:  "johndoe",
				FirstName: "John",
				LastName:  "Doe",
			},
			mention: true,
			want:    "@johndoe",
		},
		{
			name: "First and last name only",
			user: &tgbotapi.User{
				UserName:  "",
				FirstName: "John",
				LastName:  "Doe",
			},
			mention: false,
			want:    "John Doe",
		},
		{
			name: "First name only",
			user: &tgbotapi.User{
				UserName:  "",
				FirstName: "John",
				LastName:  "",
			},
			mention: false,
			want:    "John ",
		},
		{
			name: "Last name only",
			user: &tgbotapi.User{
				UserName:  "",
				FirstName: "",
				LastName:  "Doe",
			},
			mention: false,
			want:    "Doe",
		},
		{
			name: "No name info",
			user: &tgbotapi.User{
				UserName:  "",
				FirstName: "",
				LastName:  "",
			},
			mention: false,
			want:    "Хер_знает_кто_такой",
		},
		{
			name: "Username with mention overrides first name",
			user: &tgbotapi.User{
				UserName:  "testuser",
				FirstName: "Test",
				LastName:  "User",
			},
			mention: true,
			want:    "@testuser",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetUserName(tt.user, tt.mention); got != tt.want {
				t.Errorf("GetUserName() = %v, want %v", got, tt.want)
			}
		})
	}
}
