package generator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandomLocation(t *testing.T) {
	tests := []struct {
		name      string
		want      func(*testing.T, interface{})
		assertion assert.ErrorAssertionFunc
	}{
		{
			name:      "Test that result is not empty string",
			want:      func(t *testing.T, r interface{}) { assert.NotEmpty(t, r) },
			assertion: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RandomLocation()
			tt.assertion(t, err)
			tt.want(t, got)
		})
	}
}

func TestRandomChaplin(t *testing.T) {
	tests := []struct {
		name      string
		want      func(*testing.T, interface{})
		assertion assert.ErrorAssertionFunc
	}{
		{
			name:      "Test that result is not empty string",
			want:      func(t *testing.T, r interface{}) { assert.NotEmpty(t, r) },
			assertion: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RandomChaplin()
			tt.assertion(t, err)
			tt.want(t, got)
		})
	}
}
