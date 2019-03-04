package utils_test

import (
	"testing"

	"github.com/seadiaz/katalog/src/utils"
	assert "gopkg.in/go-playground/assert.v1"
)

func TestShouldFailAuthWhenPasswordNotMatch(t *testing.T) {
	input := "dummy text"

	output := utils.Serialize(input)

	assert.Equal(t, "alksdjfbakljsdfbasdkljfb", output)
}
