package main

import (
	"snippetbox.whendeadline.net/internal/assert"
	"testing"
	"time"
)

func TestHumanDate(t *testing.T) {
	tests := []struct {
		name string
		time time.Time
		want string
	}{
		{
			name: "UTC",
			time: time.Date(2022, 11, 24, 4, 20, 0, 0, time.UTC),
			want: "24 Nov 2022 at 04:20",
		},
		{
			name: "Empty",
			time: time.Time{},
			want: "",
		},
		{
			name: "UTC+3",
			time: time.Date(2022, 11, 24, 4, 20, 0, 0, time.FixedZone("Moscow", 3*60*60)),
			want: "24 Nov 2022 at 01:20",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hd := humanDate(tt.time)
			assert.Equal(t, hd, tt.want)
		})
	}

}
