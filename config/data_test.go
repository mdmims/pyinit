package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetEmbeds(t *testing.T) {
	for file, _ := range FileMap {
		data, err := GetEmbeds(file)

		assert.NoError(t, err)
		assert.NotEmpty(t, data)
	}
}
