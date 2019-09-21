package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestTree -- creates a tree and checks if the expected indexes get returned
// (indexes are assigned in the order of the nodes' creation)
func TestTree(t *testing.T) {
	var tree RWTree

	assert.Equal(t, 3, tree.GetIndex("_faster", "key", "info.json")) // 1 2 3
	assert.Equal(t, 5, tree.GetIndex("http", "GET /robots.txt"))     // 4 5
	assert.Equal(t, 6, tree.GetIndex("http", "GET /favicon.ico"))    // 4 6
	assert.Equal(t, 7, tree.GetIndex("_faster", "key", "foobar"))    // 1 2 7
	assert.Equal(t, 8, tree.GetIndex("https"))

	assert.Equal(t, 0, tree.GetIndex())
	assert.Equal(t, 1, tree.GetIndex("_faster"))
	assert.Equal(t, 2, tree.GetIndex("_faster", "key"))
	assert.Equal(t, 4, tree.GetIndex("http"))
}
