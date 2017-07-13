package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateVersion(t *testing.T) {
	assert.Equal(t, "0.1.0", createVersion("0.1.0", ""))
	assert.Equal(t, "0.1.0-abc", createVersion("0.1.0", "abc"))
  assert.Equal(t, "0.1.0-12345678", createVersion("0.1.0", "12345678"))
  assert.Equal(t, "0.1.0-12345678", createVersion("0.1.0", "123456789"))
}
