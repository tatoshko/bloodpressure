package handlerLog

import (
	"testing"
)

func TestCheck(t *testing.T) {
	tests := []struct {
		name    string
		message string
		want    bool
	}{
		{
			name:    "Valid format",
			message: "120/80/60",
			want:    true,
		},
		{
			name:    "Valid format - high values",
			message: "180/120/100",
			want:    true,
		},
		{
			name:    "Valid format - low values",
			message: "90/60/50",
			want:    true,
		},
		{
			name:    "Invalid - missing pulse",
			message: "120/80",
			want:    false,
		},
		{
			name:    "Invalid - wrong separator",
			message: "120.80.60",
			want:    false,
		},
		{
			name:    "Invalid - spaces",
			message: "120 80 60",
			want:    false,
		},
		{
			name:    "Invalid - extra text",
			message: "Test 120/80/60",
			want:    false,
		},
		{
			name:    "Invalid - single digit",
			message: "12/8/6",
			want:    false,
		},
		{
			name:    "Invalid - four digits",
			message: "1200/800/600",
			want:    false,
		},
		{
			name:    "Empty string",
			message: "",
			want:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Check(tt.message); got != tt.want {
				t.Errorf("Check() = %v, want %v", got, tt.want)
			}
		})
	}
}
