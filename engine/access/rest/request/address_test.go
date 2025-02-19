package request

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddress_InvalidParse(t *testing.T) {
	inputs := []string{
		"0x1",
		"",
		"foo",
		"1",
		"@",
		"ead892083b3e2c61222",
	}

	for _, input := range inputs {
		_, err := ParseAddress(input)
		assert.EqualError(t, err, "invalid address")
	}
}

func TestAddress_ValidParse(t *testing.T) {
	inputs := []string{
		"f8d6e0586b0a20c7",
		"f3ad66eea58c97d2",
		"0xead892083b3e2c6c",
	}

	for _, input := range inputs {
		address, err := ParseAddress(input)
		assert.NoError(t, err)
		assert.Equal(t, strings.ReplaceAll(input, "0x", ""), address.String())
	}
}
