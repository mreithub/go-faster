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

	assert.True(t, tree.Exists())
	assert.True(t, tree.Exists("_faster", "key", "foobar"))
	assert.False(t, tree.Exists("_faster", "key", "value"))

	assert.Equal(t, tree.root.GetIndex("_faster"), 1)
	assert.Equal(t, tree.root.GetIndex("_faster", "key"), 2)
	assert.Equal(t, tree.root.GetIndex("_faster", "key", "foobar"), 7)
	assert.Equal(t, tree.root.GetIndex("_faster", "key", "foobar", "bak"), -1)
}

func TestLimit(t *testing.T) {
	var tree = RWTree{Limit: 5}

	assert.Equal(t, 2, tree.GetIndex("_faster", "key"))           // 1 2
	assert.Equal(t, 4, tree.GetIndex("http", "GET /robots.txt"))  // 3 4
	assert.Equal(t, 5, tree.GetIndex("http", "GET /favicon.ico")) // 3 5(overflow)
	assert.Equal(t, 5, tree.GetIndex("_faster", "key", "foobar")) // 1 2 5(overflow)
	assert.Equal(t, 5, tree.GetIndex("https"))                    // 5(overflow)

	assert.Contains(t, tree.root.children, "_faster")
	assert.Contains(t, tree.root.children, "http")
	assert.Contains(t, tree.root.children, "_overflow")
	assert.Equal(t, 3, len(tree.root.children))

	assert.True(t, tree.Exists("_faster", "key"))
	assert.False(t, tree.Exists("_faster", "key", "foobar"))
	assert.True(t, tree.Exists("http", "GET /robots.txt"))
	assert.False(t, tree.Exists("http", "GET /favicon.ico"))
	assert.False(t, tree.Exists("https"))
}
