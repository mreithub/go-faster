package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetChild(t *testing.T) {
	var d Data

	assert.Equal(t, &d, d.GetChild()) // empty parameters -> return itself

	// check reuse
	var dbData = d.GetChild("backend", "db")
	assert.Equal(t, dbData, d.GetChild("backend", "db"))
	assert.True(t, dbData != d.GetChild("backend", "notdb"))

	// check .Children
	assert.Equal(t, 1, len(d.Children))
	assert.Contains(t, d.Children, "backend")
	assert.NotContains(t, d.Children, "db")

	var backendData = d.Children["backend"]
	assert.Equal(t, backendData, d.GetChild("backend")) // is supposed to return the existing instance
	assert.Contains(t, backendData.Children, "db")
	assert.Equal(t, backendData.Children["db"], dbData)
}
