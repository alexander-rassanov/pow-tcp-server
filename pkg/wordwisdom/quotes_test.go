package wordwisdom

import (
	"math/rand"
	"testing"
)

func TestGetRandQuote(t *testing.T) {
	rand.Seed(1)
	tests := []struct {
		name string
		want string
	}{
		{
			"Get rand quote",
			"Others see in the word of wisdom a teaching function.",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetRandQuote(); got != tt.want {
				t.Errorf("GetRandQuote() = %v, want %v", got, tt.want)
			}
		})
	}
}
