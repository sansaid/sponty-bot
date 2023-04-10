package generator

import (
	_ "embed"
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
			name: "Test that result is not an empty object",
			want: func(t *testing.T, r interface{}) {
				assert.NotNil(t, r.(Location).Name)
				assert.NotNil(t, r.(Location).Location)
			},
			assertion: assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RandomLocation("park")
			tt.assertion(t, err)
			tt.want(t, got)
		})

		t.Run(tt.name, func(t *testing.T) {
			got, err := RandomLocation("pub")
			tt.assertion(t, err)
			tt.want(t, got)
		})
	}

	t.Run("An unrecognised location returns an error", func(t *testing.T) {
		_, err := RandomLocation("doesnotexist")
		assert.Error(t, err, "location type unrecognised")
	})
}

func TestRandomPerk(t *testing.T) {
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
			got, err := RandomPerk("park")
			tt.assertion(t, err)
			tt.want(t, got)
		})

		t.Run(tt.name, func(t *testing.T) {
			got, err := RandomPerk("pub")
			tt.assertion(t, err)
			tt.want(t, got)
		})
	}
	t.Run("An unrecognised location returns an error", func(t *testing.T) {
		_, err := RandomPerk("doesnotexist")
		assert.Error(t, err, "location type unrecognised")
	})
}

func TestRandomIntro(t *testing.T) {
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
			got, err := RandomIntro()
			tt.assertion(t, err)
			tt.want(t, got)
		})
	}
}
