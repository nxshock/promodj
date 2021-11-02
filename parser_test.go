package main

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConstructURL(t *testing.T) {
	tests := []struct {
		params   url.Values
		expected string
	}{
		{nil, "https://promodj.com/music/testGenre"},
		{url.Values{"download": []string{"1"}}, "https://promodj.com/music/testGenre?download=1"},
		{url.Values{"download": []string{"1"}, "page": []string{"1"}}, "https://promodj.com/music/testGenre?download=1&page=1"},
	}

	for _, test := range tests {
		assert.Equal(t, test.expected, constructUrl("testGenre", test.params))
	}
}
