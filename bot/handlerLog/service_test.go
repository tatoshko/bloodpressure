package handlerLog

import (
	"testing"
	"time"
)

func TestLogRecord_Score(t *testing.T) {
	tests := []struct {
		name string
		lr   LogRecord
		want int
	}{
		{
			name: "Normal values",
			lr:   LogRecord{Up: 120, Down: 80},
			want: 200,
		},
		{
			name: "Zero values",
			lr:   LogRecord{Up: 0, Down: 0},
			want: 0,
		},
		{
			name: "High values",
			lr:   LogRecord{Up: 180, Down: 120},
			want: 300,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.lr.Score(); got != tt.want {
				t.Errorf("LogRecord.Score() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLogService_ComputePressureMedian(t *testing.T) {
	ls := &LogService{}
	now := time.Now()

	tests := []struct {
		name     string
		records  []*LogRecord
		wantUp   int
		wantDown int
	}{
		{
			name:     "Empty records",
			records:  nil,
			wantUp:   0,
			wantDown: 0,
		},
		{
			name: "Single record",
			records: []*LogRecord{
				{Up: 120, Down: 80, Pulse: 60, CreatedAt: now},
			},
			wantUp:   120,
			wantDown: 80,
		},
		{
			name: "Two records - first higher",
			records: []*LogRecord{
				{Up: 140, Down: 90, Pulse: 70, CreatedAt: now},
				{Up: 120, Down: 80, Pulse: 60, CreatedAt: now},
			},
			wantUp:   140,
			wantDown: 90,
		},
		{
			name: "Two records - second higher",
			records: []*LogRecord{
				{Up: 120, Down: 80, Pulse: 60, CreatedAt: now},
				{Up: 140, Down: 90, Pulse: 70, CreatedAt: now},
			},
			wantUp:   140,
			wantDown: 90,
		},
		{
			name: "Odd number of records",
			records: []*LogRecord{
				{Up: 100, Down: 60, Pulse: 50, CreatedAt: now},
				{Up: 120, Down: 80, Pulse: 60, CreatedAt: now},
				{Up: 140, Down: 90, Pulse: 70, CreatedAt: now},
			},
			wantUp:   120,
			wantDown: 80,
		},
		{
			name: "Even number - different scores, tiebreak by up",
			records: []*LogRecord{
				{Up: 120, Down: 80, Pulse: 60, CreatedAt: now},
				{Up: 130, Down: 70, Pulse: 60, CreatedAt: now},
				{Up: 110, Down: 90, Pulse: 55, CreatedAt: now},
				{Up: 100, Down: 95, Pulse: 50, CreatedAt: now},
			},
			wantUp:   130,
			wantDown: 70,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ls.ComputePressureMedian(tt.records)
			if tt.records == nil || len(tt.records) == 0 {
				if got != nil {
					t.Errorf("ComputePressureMedian() = %v, want nil", got)
				}
				return
			}
			if got.Up != tt.wantUp || got.Down != tt.wantDown {
				t.Errorf("ComputePressureMedian() Up=%v Down=%v, want Up=%v Down=%v",
					got.Up, got.Down, tt.wantUp, tt.wantDown)
			}
		})
	}
}

func TestLogService_ComputePulseMedian(t *testing.T) {
	ls := &LogService{}
	now := time.Now()

	tests := []struct {
		name      string
		records   []*LogRecord
		wantUp    int
		wantPulse int
	}{
		{
			name:      "Empty records",
			records:   nil,
			wantUp:    0,
			wantPulse: 0,
		},
		{
			name: "Single record",
			records: []*LogRecord{
				{Up: 120, Down: 80, Pulse: 70, CreatedAt: now},
			},
			wantUp:    120,
			wantPulse: 70,
		},
		{
			name: "Two records - first higher pulse",
			records: []*LogRecord{
				{Up: 120, Down: 80, Pulse: 80, CreatedAt: now},
				{Up: 120, Down: 80, Pulse: 60, CreatedAt: now},
			},
			wantUp:    120,
			wantPulse: 80,
		},
		{
			name: "Two records - second higher pulse",
			records: []*LogRecord{
				{Up: 120, Down: 80, Pulse: 60, CreatedAt: now},
				{Up: 120, Down: 80, Pulse: 80, CreatedAt: now},
			},
			wantUp:    120,
			wantPulse: 80,
		},
		{
			name: "Odd number of records",
			records: []*LogRecord{
				{Up: 100, Down: 60, Pulse: 50, CreatedAt: now},
				{Up: 120, Down: 80, Pulse: 70, CreatedAt: now},
				{Up: 140, Down: 90, Pulse: 90, CreatedAt: now},
			},
			wantUp:    120,
			wantPulse: 70,
		},
		{
			name: "Even number - different pulse tiebreak by up",
			records: []*LogRecord{
				{Up: 100, Down: 60, Pulse: 60, CreatedAt: now},
				{Up: 120, Down: 80, Pulse: 60, CreatedAt: now},
				{Up: 110, Down: 70, Pulse: 55, CreatedAt: now},
				{Up: 130, Down: 85, Pulse: 55, CreatedAt: now},
			},
			wantUp:    120,
			wantPulse: 60,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ls.ComputePulseMedian(tt.records)
			if tt.records == nil || len(tt.records) == 0 {
				if got != nil {
					t.Errorf("ComputePulseMedian() = %v, want nil", got)
				}
				return
			}
			if got.Pulse != tt.wantPulse || got.Up != tt.wantUp {
				t.Errorf("ComputePulseMedian() Up=%v Pulse=%v, want Up=%v Pulse=%v",
					got.Up, got.Pulse, tt.wantUp, tt.wantPulse)
			}
		})
	}
}
