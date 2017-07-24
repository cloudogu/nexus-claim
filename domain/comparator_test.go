package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsEqual(t *testing.T) {
	propertiesOne := make(map[string]interface{})
	propertiesOne["one"] = "1"
	sliceOne := make([]interface{}, 1)
	sliceOne[0] = propertiesOne

	propertiesTwo := make(map[string]interface{})
	propertiesTwo["one"] = "1"
	propertiesTwo["two"] = "2"
	sliceTwo := make([]interface{}, 1)
	sliceTwo[0] = propertiesTwo

	assert.True(t, IsEqual(sliceOne, sliceTwo))

	propertiesTwo["one"] = "2"
	assert.False(t, IsEqual(sliceOne, sliceTwo))
}

func TestIsEqualWithBiggerLeftSlice(t *testing.T) {
	propertiesOne := make(map[string]interface{})
	propertiesOne["one"] = "1"
	propertiesOne["two"] = "2"
	sliceOne := make([]interface{}, 2)
	sliceOne[0] = propertiesOne
	sliceOne[1] = propertiesOne

	propertiesTwo := make(map[string]interface{})
	propertiesTwo["one"] = "1"
	propertiesTwo["two"] = "2"
	sliceTwo := make([]interface{}, 1)
	sliceTwo[0] = propertiesTwo

	assert.False(t, IsEqual(sliceOne, sliceTwo))
}
