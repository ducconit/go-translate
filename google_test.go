package gotranslate

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGoogleTranslate(t *testing.T) {
	assert := assert.New(t)

	expects := map[string]map[string]bool{
		"hello": {
			"xin chào": true,
			"Xin chào": false,
			"xin Chào": false,
			"":         false,
		},
		" ": {
			" ": true,
			"":  false,
		},
	}
	for orig, value := range expects {
		client := NewGoogleTranslate()
		target, err := client.Translate(orig)
		assert.NoError(err)
		for expect, boolean := range value {
			if boolean {
				assert.Equal(target, expect)
				continue
			}
			assert.NotEqual(target, expect)
		}
	}
}
